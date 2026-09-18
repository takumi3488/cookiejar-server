FROM golang:1.27.1-alpine@sha256:4cb7ac979db5fcc41cae44b2227ba5ab8a51e8807f40d9ba4dee20a0ad960b5b AS builder
WORKDIR /usr/src/app
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN go build -o /usr/local/bin/writer ./cmd/writer
RUN go build -o /usr/local/bin/reader ./cmd/reader


FROM gcr.io/distroless/static-debian12:nonroot@sha256:afa5c872c891853ca7fcf1f12c3edb23f7eeef36189728842dd51042ff57f7ab AS writer
COPY --from=builder /usr/local/bin/writer /app
ENTRYPOINT ["/app"]

FROM gcr.io/distroless/static-debian12:nonroot@sha256:afa5c872c891853ca7fcf1f12c3edb23f7eeef36189728842dd51042ff57f7ab AS reader
COPY --from=builder /usr/local/bin/reader /app
ENTRYPOINT ["/app"]
