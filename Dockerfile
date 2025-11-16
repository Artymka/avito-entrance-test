FROM golang:1.25-alpine AS builder

WORKDIR /builder

COPY go.mod go.sum ./

RUN go mod download

COPY ./internal/ ./internal/
COPY ./cmd/server/ ./cmd/server/

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags='-w -s' -o server ./cmd/server/main.go

FROM alpine:latest AS runner

WORKDIR /app

COPY --from=builder /builder/server .
COPY ./config/deploy.yaml ./config/local.yaml

EXPOSE 8080

CMD ["./server"]