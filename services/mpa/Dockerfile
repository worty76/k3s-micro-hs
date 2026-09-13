FROM golang:1.27.0-alpine AS builder

WORKDIR /app

COPY libs ./libs

WORKDIR /app/services/mpa
COPY services/mpa/go.mod services/mpa/go.sum ./
RUN go mod download

COPY services/mpa ./

RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /mpa ./cmd/server

FROM alpine:3.22

RUN addgroup -S app && adduser -S app -G app

WORKDIR /app

COPY --from=builder /mpa ./mpa

USER app

ENTRYPOINT ["./mpa"]
