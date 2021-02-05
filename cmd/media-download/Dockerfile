FROM golang:alpine AS build

WORKDIR /go/src/github.com/reecerussell/open-social

COPY client/ client/
COPY core/ core/
COPY util/ util/
COPY cmd/media-download cmd/media-download/
COPY go.mod go.sum ./

RUN go mod download

WORKDIR /go/src/github.com/reecerussell/open-social/cmd/media-download

RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o /app/main main.go

FROM scratch

COPY --from=build /app/main main

CMD ["./main"]