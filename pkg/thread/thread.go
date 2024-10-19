package thread

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/pkg/errors"

	. "github.com/helloacai/spindle/pkg/util" // Hex
)

var threadMap map[string]*Thread // uid -> thread
var lock sync.Mutex
var listenerMap map[string]map[string]*Listener // uid -> request id -> Listener

func init() {
	// TODO: this should really be persisted somewhere and not just sitting in memory
	threadMap = map[string]*Thread{}
	listenerMap = map[string]map[string]*Listener{}
}

type EntryType string

const (
	EntryType_Request  EntryType = "request"
	EntryType_Waiting  EntryType = "waiting"
	EntryType_Complete EntryType = "complete"
)

type HexByteSlice []byte

func (s HexByteSlice) MarshalJSON() ([]byte, error) {
	return json.Marshal(Hex(s))
}

type Entry struct {
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
}

func (t *Thread) notify() {
	toRemove := []string{}
	for requestID, listener := range listenerMap[Hex(t.UID)] {
		if errors.Is(listener.ctx.Err(), context.Canceled) {
			// remove unused listeners
			close(listener.ch)
			toRemove = append(toRemove, requestID)
			continue
		}
		listener.ch <- t.Context[len(t.Context)-1]
	}
	for _, requestID := range toRemove {
		delete(listenerMap[Hex(t.UID)], requestID)
	}
}

func (t *Thread) append(typ EntryType, originator []byte, message string) *Thread {
	t.Context = append(t.Context, &Entry{
		Type:       typ,
		Originator: originator,
		Message:    message,
	})
	return t
}

func (t *Thread) Append(typ EntryType, originator []byte, message string) {
	lock.Lock()
	defer lock.Unlock()

	t.append(typ, originator, message)
	t.notify()
}

func Request(uid []byte, parentUID []byte, aciUID []byte, requester []byte, requestRef string) *Thread {
	lock.Lock()
	defer lock.Unlock()

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
					Type:       EntryType_Request,
					Originator: requester,
					Message:    requestRef, // TODO: ref -> actual request
				},
			},
		}
		threadMap[Hex(uid)] = t
	} else {
		// update case
		t.append(EntryType_Request, requester, requestRef)
	}

	// notify listeners
	t.notify()

	return t
}

type Listener struct {
	ctx context.Context
	ch  chan *Entry
}

func Listen(ctx context.Context, uid []byte, requestID string) (<-chan *Entry, error) {
	lock.Lock()
	defer lock.Unlock()
	if _, exists := threadMap[Hex(uid)]; !exists {
		return nil, errors.New("thread does not exist")
	}

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

	go func() {
		lock.Lock()
		defer lock.Unlock()
		// throw all existing entries onto the channel
		for _, entry := range threadMap[Hex(uid)].Context {
			l.ch <- entry
		}
	}()

	return l.ch, nil
}
