.PHONY: proto
proto:
	protoc --go_out=. --go_opt=Msubstreams/proto/contract.proto=pb/contract/v1 substreams/proto/contract.proto
