FROM golang:1.20-alpine AS builder
WORKDIR /app
COPY go.mod .
COPY cmd ./cmd
RUN go build -o sre-demo ./cmd/app

FROM alpine:3.18
WORKDIR /
COPY --from=builder /app/sre-demo /sre-demo
EXPOSE 8080
ENTRYPOINT ["/sre-demo"]
