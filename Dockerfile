FROM golang:1.25.5-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s" \
    -o /sidecar \
    ./cmd/proxy/sidecar.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

COPY --from=builder /sidecar .

ENV SIDECAR_PORT=8000
EXPOSE 8000

ENTRYPOINT ["./sidecar"]