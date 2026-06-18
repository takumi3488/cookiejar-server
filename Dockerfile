FROM golang:1.26.4-alpine@sha256:3ad57304ad93bbec8548a0437ad9e06a455660655d9af011d58b993f6f615648 AS builder
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
