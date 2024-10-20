package manager

import (
	"context"

	"github.com/pkg/errors"

	v1 "github.com/helloacai/spindle/pb/contract/v1"
	"github.com/helloacai/spindle/pkg/aciregistry"
	"github.com/helloacai/spindle/pkg/agent"
	"github.com/helloacai/spindle/pkg/log"
	"github.com/helloacai/spindle/pkg/thread"
	. "github.com/helloacai/spindle/pkg/util" // Hex
)

func HandleRequest(_ context.Context, call *v1.Acs_RequestCall) error {
	log.Debug().
		Uint64("block", call.CallBlockNumber).
		Str("aci_uid", Hex(call.AciUid)).
		Str("parent_thread_uid", Hex(call.ParentThreadUid)).
		Str("thread_uid", Hex(call.ThreadUid)).
		Str("request_ref", call.RequestRef).
		Msg("[call] request")
	return nil
}

func HandleRequested(ctx context.Context, event *v1.Acs_Requested) error {
	log.Debug().
		Uint64("block", event.EvtBlockNumber).
		Str("aci_uid", Hex(event.AciUid)).
		Str("thread_uid", Hex(event.ThreadUid)).
		Str("parent_thread_uid", Hex(event.ParentThreadUid)).
		Str("requester", Hex(event.Requester)).
		Str("request_ref", event.RequestRef).
		Msg("[event] requested")

	// log the request in the thread
	t, isNew := thread.Request(event.ThreadUid, event.ParentThreadUid, event.AciUid, event.Requester, event.RequestRef)

	// fetch the aci metadata
	if t.ACIMetadata == nil {
		aciMetadata, err := aciregistry.Get(ctx, event.AciUid)
		if err != nil {
			return errors.Wrap(err, "error fetching aci registry metadata")
		}
		t.ACIMetadata = aciMetadata
	}
	log.Debug().Object("aci_metadata", t.ACIMetadata).Msg("fetched aci metadata")

	// call the agent API
	agentResponse, err := agent.Call(ctx, t.ACIMetadata, event.RequestRef, event.ParentThreadUid, event.ThreadUid, isNew)
	if err != nil {
		return errors.Wrap(err, "error calling agent")
	}
	log.Debug().Object("agent_response", agentResponse).Msg("agent responded")

	// log the result in the thread
	for _, msg := range agentResponse.Messages {
		switch msg.Status {
		case agent.Status_Debug:
			t.Append(thread.EntryType_Debug, event.AciUid, msg.Message)
		case agent.Status_Info:
			t.Append(thread.EntryType_Info, event.AciUid, msg.Message)
		case agent.Status_Waiting:
			t.Append(thread.EntryType_Waiting, event.AciUid, msg.Message)
		case agent.Status_Complete:
			t.Append(thread.EntryType_Complete, event.AciUid, msg.Message)
		default:
			return errors.New("invalid agent response status: " + string(msg.Status))
		}
	}

	return nil
}

func HandleThreadFunded(_ context.Context, event *v1.Acs_ThreadFunded) error {
	log.Debug().
		Uint64("block", event.EvtBlockNumber).
		Str("thread_uid", Hex(event.ThreadUid)).
		Str("funder", Hex(event.Funder)).
		Uint64("funding_ammount", event.FundingAmount).
		Msg("[event] thread funded")
	return nil
}

func Handle(ctx context.Context, eventsCalls *v1.EventsCalls) error {
	for _, call := range eventsCalls.Calls.AcsCallRequests {
		if call == nil {
			continue
		}
		if err := HandleRequest(ctx, call); err != nil {
			return errors.Wrap(err, "error handling request")
		}
	}

	for _, event := range eventsCalls.Events.AcsRequesteds {
		if event == nil {
			continue
		}
		if err := HandleRequested(ctx, event); err != nil {
			return errors.Wrap(err, "error handling requested")
		}
	}

	for _, event := range eventsCalls.Events.AcsThreadFundeds {
		if event == nil {
			continue
		}
		if err := HandleThreadFunded(ctx, event); err != nil {
			return errors.Wrap(err, "error handling thread funded")
		}
	}

	return nil
}
