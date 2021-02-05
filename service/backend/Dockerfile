FROM golang:alpine AS build

WORKDIR /go/src/github.com/reecerussell/open-social

COPY client/ client/
COPY core/ core/
COPY util/ util/
COPY service/backend service/backend/
COPY go.mod go.sum ./

RUN go mod download

WORKDIR /go/src/github.com/reecerussell/open-social/service/backend

RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o /app/main main.go

FROM scratch

COPY --from=build /app/main main

CMD ["./main"]