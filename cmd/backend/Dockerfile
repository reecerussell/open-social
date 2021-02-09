FROM golang:alpine AS build

ENV USER=app
ENV UID=10001

RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistant" \
    --shell "/sbin/nologin" \
    --no-create-home \\
    --uid "${UID}" \
    "${USER}"

WORKDIR /go/src/github.com/reecerussell/open-social

COPY . .

RUN go mod download

WORKDIR /go/src/github.com/reecerussell/open-social/cmd/backend

RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-w -s" -o /app/main main.go

FROM scratch

COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group
COPY --from=build /app/main main

USER app:app

CMD ["./main"]