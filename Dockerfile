FROM golang:1.24.4 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .


RUN CGO_ENABLED=0 GOOS=linux go build -o /go-service ./cmd/main.go


FROM alpine:3.18
RUN apk add --no-cache netcat-openbsd

WORKDIR /app

COPY --from=builder /go-service .
COPY --from=builder /app/configs ./configs
COPY --from=builder /app/web ./web

ENV GIN_MODE=release

EXPOSE 8081

COPY entrypoint.sh .
RUN chmod +x ./entrypoint.sh

CMD ["./entrypoint.sh"]