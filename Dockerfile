FROM golang:1.24-bookworm AS build-env

WORKDIR /app
COPY . .

RUN go mod download

RUN CGO_ENABLED=1 go build -o MT_PowerUsage

CMD ["./MT_PowerUsage"]
