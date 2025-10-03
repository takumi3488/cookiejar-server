ARG GO_VERSION=1.25.1
FROM golang:${GO_VERSION}-alpine AS builder

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o /usr/local/bin/app ./cmd


FROM gcr.io/distroless/static-debian12

COPY --from=builder /usr/local/bin/app /app

ENTRYPOINT ["/app"]
