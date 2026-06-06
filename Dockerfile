FROM golang:1.26.4-alpine@sha256:f23e8b227fb4493eabe03bede4d5a32d04092da71962f1fb79b5f7d1e6c2a17f AS builder
WORKDIR /usr/src/app
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN go build -o /usr/local/bin/writer ./cmd/writer
RUN go build -o /usr/local/bin/reader ./cmd/reader


FROM gcr.io/distroless/static-debian12:nonroot@sha256:d093aa3e30dbadd3efe1310db061a14da60299baff8450a17fe0ccc514a16639 AS writer
COPY --from=builder /usr/local/bin/writer /app
ENTRYPOINT ["/app"]

FROM gcr.io/distroless/static-debian12:nonroot@sha256:d093aa3e30dbadd3efe1310db061a14da60299baff8450a17fe0ccc514a16639 AS reader
COPY --from=builder /usr/local/bin/reader /app
ENTRYPOINT ["/app"]
