FROM golang:alpine AS build

RUN apk update && apk add --no-cache ca-certificates

WORKDIR /go/src/github.com/reecerussell/open-social

COPY . .

RUN go mod download

WORKDIR /go/src/github.com/reecerussell/open-social/cmd/media

RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o /app/main main.go

FROM scratch

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /app/main main

CMD ["./main"]