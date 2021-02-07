FROM golang:alpine AS build

WORKDIR /go/src/github.com/reecerussell/open-social

COPY . .

RUN go mod download

WORKDIR /go/src/github.com/reecerussell/open-social/service/posts

RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o /app/main main.go

FROM scratch

COPY --from=build /app/main main

CMD ["./main"]