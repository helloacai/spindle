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
	EntryType_Request     EntryType = "request"
	EntryType_Update      EntryType = "update"
	EntryType_Info        EntryType = "info"
	EntryType_Debug       EntryType = "debug"
	EntryType_Waiting     EntryType = "waiting"
	EntryType_Complete    EntryType = "complete"
	EntryType_Subcomplete EntryType = "subcomplete"
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
		Originator: getOriginator(originator),
		Message:    message,
	})

	if t.ParentUID != nil {
		parent, exists := threadMap[Hex(t.ParentUID)]
		if !exists {
			log.Err(errors.New("parent " + Hex(t.ParentUID) + " does not exist")).Msg("parent does not exist")
		} else {
			if typ != EntryType_Complete {
				parent.append(typ, originator, message)
			} else {
				parent.append(EntryType_Subcomplete, originator, message)
			}
		}
	}

	return t
}

// HACK: not enough time to pull this from aci registry
func getOriginator(uid []byte) []byte {
	uidHex := Hex(uid)
	var originator []byte
	var err error
	if uidHex == "0xaf3a33c2f95a9e41e54c386aad4d260b2f3fe73a73353d3c439871bbf2301e41" {
		originator, err = FromHex("0x980b3a841817e0c7cbfc921bfccf796fa84c80f3")
	} else if uidHex == "0xa3f374f49528ef97f2c6adad7931e87373782d5e9b965de2e61554828275a033" {
		originator, err = FromHex("0x7b5d5f79934f4995a5d9214b6701e98e903d26c7")
	} else if uidHex == "0x822434c25a9837f0e7244090c1558663dee097f16f7623f0bf461c8afee4c55b" {
		originator, err = FromHex("0x17f5e274c2c8aa9471e5320cb68df543db90d083")
	} else if uidHex == "0xeaac656a5054ef4a92f34da8870b97a7a3037f20181ad169956b4a631903d466" {
		originator, err = FromHex("0xc04c0b99b0b50e61f6a71f43e18871803c2ba2af")
	} else {
		originator = uid
	}

	if err != nil {
		log.Err(err).Msg("error getting originator")
	}

	return originator
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
					Originator: getOriginator(requester),
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

func Get(ctx context.Context, uid []byte) (*Thread, error) {
	threadMapLock.RLock()
	defer threadMapLock.RUnlock()
	t, exists := threadMap[Hex(uid)]
	if !exists {
		return nil, errors.New("thread does not exist")
	}
	return t, nil
}

func GetByParentAndAci(ctx context.Context, parentUID, aciUID []byte) (*Thread, error) {
	threadMapLock.RLock()
	defer threadMapLock.RUnlock()
	for _, t := range threadMap {
		if Hex(t.ParentUID) == Hex(parentUID) && Hex(t.AciUID) == Hex(aciUID) {
			return t, nil
		}
	}
	return nil, errors.New("thread not found")
}

type Listener struct {
	ctx context.Context
	ch  chan *Entry
}

func Listen(ctx context.Context, uid []byte, requestID string) (<-chan *Entry, error) {
	threadMapLock.RLock()
	if _, exists := threadMap[Hex(uid)]; !exists {
		threadMapLock.RUnlock()
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
