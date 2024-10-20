package thread

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/pkg/errors"

	"github.com/helloacai/spindle/pkg/aciregistry"
	"github.com/helloacai/spindle/pkg/log"
	. "github.com/helloacai/spindle/pkg/util" // Hex
)

var threadMap map[string]*Thread // uid -> thread
var threadMapLock sync.RWMutex
var listenerMap map[string]map[string]*Listener // uid -> request id -> Listener
var listenerMapLock sync.Mutex

func init() {
	// TODO: this should really be persisted somewhere and not just sitting in memory
	threadMap = map[string]*Thread{}
	listenerMap = map[string]map[string]*Listener{}
}

type EntryType string

const (
	EntryType_Request  EntryType = "request"
	EntryType_Update   EntryType = "update"
	EntryType_Info     EntryType = "info"
	EntryType_Debug    EntryType = "debug"
	EntryType_Waiting  EntryType = "waiting"
	EntryType_Complete EntryType = "complete"
)

type HexByteSlice []byte

func (s HexByteSlice) MarshalJSON() ([]byte, error) {
	return json.Marshal(Hex(s))
}

type Entry struct {
	ID         int          `json:"id"`
	Type       EntryType    `json:"type"`
	Originator HexByteSlice `json:"originator"`
	Message    string       `json:"message"`
}

type Thread struct {
	UID       HexByteSlice `json:"uid"`
	ParentUID HexByteSlice `json:"parentUID"`
	AciUID    HexByteSlice `json:"aciUID"`
	Requester HexByteSlice `json:"requester"`
	Context   []*Entry     `json:"context"`

	ACIMetadata *aciregistry.Metadata `json:"-"`
}

func (t *Thread) notify() {
	toRemove := []string{}
	for requestID, listener := range listenerMap[Hex(t.UID)] {
		if errors.Is(listener.ctx.Err(), context.Canceled) {
			// remove unused listeners
			log.Debug().Str("request_id", requestID).Str("thread", Hex(t.UID)).Msg("closing listener: context canceled")
			close(listener.ch)
			toRemove = append(toRemove, requestID)
			continue
		}

		listener.ch <- t.Context[len(t.Context)-1]

		// if the entry type is complete, the stream can end
		if t.Context[len(t.Context)-1].Type == EntryType_Complete {
			log.Debug().Str("request_id", requestID).Str("thread", Hex(t.UID)).Msg("closing listener: thread complete")
			close(listener.ch)
			toRemove = append(toRemove, requestID)
			continue
		}
	}
	for _, requestID := range toRemove {
		delete(listenerMap[Hex(t.UID)], requestID)
	}

	if t.ParentUID != nil {
		parent, exists := threadMap[Hex(t.ParentUID)]
		if !exists {
			log.Err(errors.New("parent " + Hex(t.ParentUID) + " does not exist")).Msg("parent does not exist")
		} else {
			parent.notify()
		}
	}
}

func (t *Thread) append(typ EntryType, originator []byte, message string) *Thread {
	t.Context = append(t.Context, &Entry{
		ID:         len(t.Context),
		Type:       typ,
		Originator: originator,
		Message:    message,
	})

	if t.ParentUID != nil {
		parent, exists := threadMap[Hex(t.ParentUID)]
		if !exists {
			log.Err(errors.New("parent " + Hex(t.ParentUID) + " does not exist")).Msg("parent does not exist")
		} else {
			parent.append(typ, originator, message)
		}
	}

	return t
}

func (t *Thread) Append(typ EntryType, originator []byte, message string) {
	threadMapLock.Lock()
	defer threadMapLock.Unlock()

	t.append(typ, originator, message)

	listenerMapLock.Lock()
	t.notify()
	listenerMapLock.Unlock()
}

func Request(uid []byte, parentUID []byte, aciUID []byte, requester []byte, requestRef string) (*Thread, bool) {
	threadMapLock.Lock()
	defer threadMapLock.Unlock()

	isNew := true
	t, exists := threadMap[Hex(uid)]
	if !exists {
		// new case
		t = &Thread{
			UID:       uid,
			ParentUID: parentUID,
			AciUID:    aciUID,
			Requester: requester,
			Context: []*Entry{
				{
					ID:         0,
					Type:       EntryType_Request,
					Originator: requester,
					Message:    requestRef, // TODO: ref -> actual request
				},
			},
		}
		threadMap[Hex(uid)] = t
	} else {
		// update case
		isNew = false
		t.append(EntryType_Request, requester, requestRef)
	}

	// notify listeners
	listenerMapLock.Lock()
	t.notify()
	listenerMapLock.Unlock()

	return t, isNew
}

type Listener struct {
	ctx context.Context
	ch  chan *Entry
}

func Listen(ctx context.Context, uid []byte, requestID string) (<-chan *Entry, error) {
	threadMapLock.RLock()
	if _, exists := threadMap[Hex(uid)]; !exists {
		return nil, errors.New("thread does not exist")
	}
	threadMapLock.RUnlock()

	listenerMapLock.Lock()
	listeners, exists := listenerMap[Hex(uid)]
	if !exists {
		listeners = map[string]*Listener{}
		listenerMap[Hex(uid)] = listeners
	}
	l := &Listener{
		ctx: ctx,
		ch:  make(chan *Entry),
	}
	listeners[requestID] = l
	listenerMapLock.Unlock()

	go func() {
		threadMapLock.RLock()
		defer threadMapLock.RUnlock()
		// throw all existing entries onto the channel
		for _, entry := range threadMap[Hex(uid)].Context {
			l.ch <- entry
		}
	}()

	return l.ch, nil
}
