.PHONY: proto
proto:
	protoc --go_out=. --go_opt=Msubstreams/proto/contract.proto=pb/contract/v1 substreams/proto/contract.proto

.PHONY: build
build:
	go build -mod=vendor -o spindle .

.PHONY: clean
clean:
	rm -f spindle

.PHONY: docker
docker-build:
	docker build . -t spindle-test

.PHONY: docker
docker-run: docker-build
	docker run --rm -v "$(shell pwd)"/secrets:/etc/secrets spindle-test:latest
