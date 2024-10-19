FROM golang:1.23.2 AS builder

WORKDIR /build
ADD . ./
RUN ls -al && go build -mod=vendor -o spindle .

FROM gcr.io/distroless/base-debian12:latest

WORKDIR /

COPY substreams/*.spkg /
COPY --from=builder /build/spindle /spindle

EXPOSE 8080

ENTRYPOINT ["/spindle"]

CMD ["sink", "amoy.substreams.pinax.network:443", "spindle-v0.1.0.spkg", "map_events_calls", "--development-mode"]
