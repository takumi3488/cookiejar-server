FROM golang:1.26.1-alpine@sha256:2389ebfa5b7f43eeafbd6be0c3700cc46690ef842ad962f6c5bd6be49ed82039 AS builder
WORKDIR /usr/src/app
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN go build -o /usr/local/bin/writer ./cmd/writer
RUN go build -o /usr/local/bin/reader ./cmd/reader


FROM gcr.io/distroless/static-debian12:nonroot@sha256:a9329520abc449e3b14d5bc3a6ffae065bdde0f02667fa10880c49b35c109fd1 AS writer
COPY --from=builder /usr/local/bin/writer /app
ENTRYPOINT ["/app"]

FROM gcr.io/distroless/static-debian12:nonroot@sha256:a9329520abc449e3b14d5bc3a6ffae065bdde0f02667fa10880c49b35c109fd1 AS reader
COPY --from=builder /usr/local/bin/reader /app
ENTRYPOINT ["/app"]
