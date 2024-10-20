FROM golang:1.23.2 AS builder

WORKDIR /build
ADD . ./
RUN go build -mod=vendor -o spindle .

FROM gcr.io/distroless/base-debian12:latest

WORKDIR /app

COPY substreams/*.spkg ./

COPY --from=builder /build/spindle ./spindle

EXPOSE 8080

ENTRYPOINT ["./spindle", "sink", "polygon.streamingfast.io:443", "/app/spindle-v0.1.0.spkg", "map_events_calls"]
