FROM golang:1.18 as builder
RUN mkdir /build
ADD . /build/
WORKDIR /build

FROM alpine:latest
COPY --from=builder /build/main . 

ENV HTTP_RECEIVER="8080" \
    GRPC_SENDER="6006" \
    METRICS_ADDRESS="443"

CMD ["./main"]