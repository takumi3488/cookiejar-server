ARG GO_VERSION=1.25.1

FROM golang:${GO_VERSION}-alpine AS builder
WORKDIR /usr/src/app
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN go build -o /usr/local/bin/writer ./cmd/writer
RUN go build -o /usr/local/bin/reader ./cmd/reader


FROM gcr.io/distroless/static-debian12:nonroot AS writer
COPY --from=builder /usr/local/bin/writer /app
ENTRYPOINT ["/app"]

FROM gcr.io/distroless/static-debian12:nonroot AS reader
COPY --from=builder /usr/local/bin/reader /app
ENTRYPOINT ["/app"]
