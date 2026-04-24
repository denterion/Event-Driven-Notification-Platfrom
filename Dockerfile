FROM golang:1.25.0 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Select which cmd/* to build (default: event-api)
ARG APP=event-api
RUN go build -o /out/app ./cmd/${APP}

FROM debian:bookworm-slim

RUN apt-get update \
  && apt-get install -y --no-install-recommends ca-certificates \
  && rm -rf /var/lib/apt/lists/*

WORKDIR /app
COPY --from=builder /out/app ./app

CMD ["./app"]