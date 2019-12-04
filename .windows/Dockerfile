FROM golang as builder
COPY main.go c:/go/src/github.com/fubarhouse/pygmy-go/
COPY go.mod c:/go/src/github.com/fubarhouse/pygmy-go/
COPY cmd/ c:/go/src/github.com/fubarhouse/pygmy-go/cmd/
COPY service/ c:/go/src/github.com/fubarhouse/pygmy-go/service/

WORKDIR c:/go/src/github.com/fubarhouse/pygmy-go/

ENV GOROOT=c:/go
ENV GO111MODULE=on
ENV GOOS=windows
ENV GOARCH=amd64

RUN go mod vendor
ENV GO111MODULE=off
RUN rm c:/go/src/github.com/fubarhouse/pygmy-go/go.mod

RUN go build -o pygmy-go.exe .
RUN mkdir c:/app
RUN cp pygmy-go.exe c:/app/pygmy-go.exe

WORKDIR /app