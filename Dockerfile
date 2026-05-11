FROM golang:1.25.0 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG APP=event-api
RUN CGO_ENABLED=0 go build -o /out/app ./cmd/${APP}

FROM scratch

WORKDIR /app
COPY --from=builder /out/app ./app
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

CMD ["./app"]
