mod abi;
mod pb;
use hex_literal::hex;
use pb::contract::v1 as contract;
use substreams::Hex;
use substreams_ethereum::pb::eth::v2 as eth;
use substreams_ethereum::Event;

#[allow(unused_imports)]
use num_traits::cast::ToPrimitive;
use std::str::FromStr;
use substreams::scalar::BigDecimal;

substreams_ethereum::init!();

const ACS_TRACKED_CONTRACT: [u8; 20] = hex!("d4cd7d0670c1197279ae1adf3507514f41971028");

fn map_acs_events(blk: &eth::Block, events: &mut contract::Events) {
    events.acs_requesteds.append(&mut blk
        .receipts()
        .flat_map(|view| {
            view.receipt.logs.iter()
                .filter(|log| log.address == ACS_TRACKED_CONTRACT)
                .filter_map(|log| {
                    if let Some(event) = abi::acs_contract::events::Requested::match_and_decode(log) {
                        return Some(contract::AcsRequested {
                            evt_tx_hash: Hex(&view.transaction.hash).to_string(),
                            evt_index: log.block_index,
                            evt_block_time: Some(blk.timestamp().to_owned()),
                            evt_block_number: blk.number,
                            aci_uid: Vec::from(event.aci_uid),
                            parent_thread_uid: Vec::from(event.parent_thread_uid),
                            request_ref: event.request_ref,
                            requester: event.requester,
                            thread_uid: Vec::from(event.thread_uid),
                        });
                    }

                    None
                })
        })
        .collect());
    events.acs_thread_fundeds.append(&mut blk
        .receipts()
        .flat_map(|view| {
            view.receipt.logs.iter()
                .filter(|log| log.address == ACS_TRACKED_CONTRACT)
                .filter_map(|log| {
                    if let Some(event) = abi::acs_contract::events::ThreadFunded::match_and_decode(log) {
                        return Some(contract::AcsThreadFunded {
                            evt_tx_hash: Hex(&view.transaction.hash).to_string(),
                            evt_index: log.block_index,
                            evt_block_time: Some(blk.timestamp().to_owned()),
                            evt_block_number: blk.number,
                            funder: event.funder,
                            funding_amount: event.funding_amount.to_u64(),
                            thread_uid: Vec::from(event.thread_uid),
                        });
                    }

                    None
                })
        })
        .collect());
}
fn map_acs_calls(blk: &eth::Block, calls: &mut contract::Calls) {
    calls.acs_call_requests.append(&mut blk
        .transactions()
        .flat_map(|tx| {
            tx.calls.iter()
                .filter(|call| call.address == ACS_TRACKED_CONTRACT && abi::acs_contract::functions::Request::match_call(call))
                .filter_map(|call| {
                    match abi::acs_contract::functions::Request::decode(call) {
                        Ok(decoded_call) => {
                            let output_param0 = match abi::acs_contract::functions::Request::output(&call.return_data) {
                                Ok(output_param0) => {output_param0}
                                Err(_) => Default::default(),
                            };
                            
                            Some(contract::AcsRequestCall {
                                call_tx_hash: Hex(&tx.hash).to_string(),
                                call_block_time: Some(blk.timestamp().to_owned()),
                                call_block_number: blk.number,
                                call_ordinal: call.begin_ordinal,
                                call_success: !call.state_reverted,
                                aci_uid: Vec::from(decoded_call.aci_uid),
                                output_param0: Vec::from(output_param0),
                                parent_thread_uid: Vec::from(decoded_call.parent_thread_uid),
                                request_ref: decoded_call.request_ref,
                                thread_uid: Vec::from(decoded_call.thread_uid),
                            })
                        },
                        Err(_) => None,
                    }
                })
        })
        .collect());
}

#[substreams::handlers::map]
fn map_events_calls(
    events: contract::Events,
    calls: contract::Calls,
) -> Result<contract::EventsCalls, substreams::errors::Error> {
    Ok(contract::EventsCalls {
        events: Some(events),
        calls: Some(calls),
    })
}
#[substreams::handlers::map]
fn map_events(blk: eth::Block) -> Result<contract::Events, substreams::errors::Error> {
    let mut events = contract::Events::default();
    map_acs_events(&blk, &mut events);
    Ok(events)
}
#[substreams::handlers::map]
fn map_calls(blk: eth::Block) -> Result<contract::Calls, substreams::errors::Error> {
let mut calls = contract::Calls::default();
    map_acs_calls(&blk, &mut calls);
    Ok(calls)
}

