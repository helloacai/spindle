package manager

import (
	"context"

	v1 "github.com/helloacai/spindle/pb/contract/v1"
	"github.com/helloacai/spindle/pkg/log"
	"github.com/helloacai/spindle/pkg/thread"
	. "github.com/helloacai/spindle/pkg/util" // Hex
)

func HandleRequest(_ context.Context, call *v1.Acs_RequestCall) error {
	log.Info().
		Uint64("block", call.CallBlockNumber).
		Str("aci_uid", Hex(call.AciUid)).
		Str("parent_thread_uid", Hex(call.ParentThreadUid)).
		Str("thread_uid", Hex(call.ThreadUid)).
		Str("request_ref", call.RequestRef).
		Msg("[call] query")
	return nil
}

func HandleRequested(_ context.Context, event *v1.Acs_Requested) error {
	log.Info().
		Uint64("block", event.EvtBlockNumber).
		Str("aci_uid", Hex(event.AciUid)).
		Str("thread_uid", Hex(event.ThreadUid)).
		Str("parent_thread_uid", Hex(event.ParentThreadUid)).
		Str("requester", Hex(event.Requester)).
		Str("request_ref", event.RequestRef).
		Msg("[event] new thread")
	thread.Request(event.ThreadUid, event.ParentThreadUid, event.AciUid, event.Requester, event.RequestRef)
	return nil
}

func HandleThreadFunded(_ context.Context, event *v1.Acs_ThreadFunded) error {
	log.Info().
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
			return err
		}
	}

	for _, event := range eventsCalls.Events.AcsRequesteds {
		if event == nil {
			continue
		}
		if err := HandleRequested(ctx, event); err != nil {
			return err
		}
	}

	for _, event := range eventsCalls.Events.AcsThreadFundeds {
		if event == nil {
			continue
		}
		if err := HandleThreadFunded(ctx, event); err != nil {
			return err
		}
	}

	return nil
}
