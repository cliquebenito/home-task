FROM golang:1.24 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o app ./cmd

FROM debian:bookworm-slim

WORKDIR /root/

COPY --from=builder /app/app .

RUN mkdir -p /app/cfg
COPY --from=builder /app/cmd/app/cfg/config.yaml /app/cfg/config.yaml

EXPOSE 8080

CMD ["./app"]