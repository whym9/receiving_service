FROM golang:1.18-alpine:latest as builder
WORKDIR /receiver
COPY . .
RUN go build -o main main.go


FROM alpine:latest
WORKDIR /receiver
COPY --from=builder /receiver/main . 

ENV HTTP_RECEIVER="localhost:8080" \
    GRPC_SENDER=":6006" \
    PROMETHEUS_ADDRESS="localhost:443"

CMD ["receiver/main"]