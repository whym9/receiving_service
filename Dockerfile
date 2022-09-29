FROM golang:1.18-alpine3.14 as builder
WORKDIR /receiver
COPY . .
RUN go build -o main main.go

FROM alpine:3.14
WORKDIR /receiver
COPY --from=builder /receiver/main .

CMD [ "/receiver/main" ]

