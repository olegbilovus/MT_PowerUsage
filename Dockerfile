FROM golang:1.24-bookworm AS build-env

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=1 go build -o MT_PowerUsage

FROM debian:bookworm-slim

WORKDIR /app

COPY --from=build-env /app/MT_PowerUsage /app/MT_PowerUsage

CMD ["./MT_PowerUsage"]
