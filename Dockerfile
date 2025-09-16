FROM golang:1.25-trixie AS build-env

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=1 go build -o /app

FROM debian:trixie-slim

COPY --from=build-env /app /app

CMD ["/app"]
