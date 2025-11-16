FROM golang:1.25-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o peerly ./cmd/peerly

FROM alpine:3.20

WORKDIR /app

COPY --from=builder /app/peerly /usr/local/bin/peerly
COPY configs ./configs

ENTRYPOINT ["/usr/local/bin/peerly"]
