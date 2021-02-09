FROM golang:alpine AS base

RUN apk update && apk add --no-cache ca-certificates tzdata && update-ca-certificates

ENV USER=app
ENV UID=10001

RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistant" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    "${USER}"

FROM base AS deps
WORKDIR /go/src/github.com/reecerussell/open-social

COPY *.* ./
COPY client/ client/
COPY util/ util/
COPY cmd/auth/ cmd/auth/

RUN go mod download
RUN go mod verify

FROM deps AS build
WORKDIR /go/src/github.com/reecerussell/open-social

ENV GOOS=linux
ENV GOARCH=amd64
ENV CGO_ENABLED=0

RUN go test ./...
RUN go build -ldflags="-w -s" -o /app/main cmd/auth/main.go

COPY cmd/auth/config.json /app/config.json

FROM scratch

COPY --from=base /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=base /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=base /etc/passwd /etc/passwd
COPY --from=base /etc/group /etc/group
COPY --from=build /app/config.json config.json
COPY --from=build /app/main main

USER ${UID}

ARG PORT=9292
EXPOSE ${PORT}
ENV PORT=${PORT}

CMD ["./main"]