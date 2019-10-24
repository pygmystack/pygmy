FROM golang as builder
COPY main.go /go/src/github.com/fubarhouse/pygmy/
COPY go.sum /go/src/github.com/fubarhouse/pygmy/
COPY go.mod /go/src/github.com/fubarhouse/pygmy/
COPY cmd/ /go/src/github.com/fubarhouse/pygmy/cmd/
COPY service/ /go/src/github.com/fubarhouse/pygmy/service/

WORKDIR /go/src/github.com/fubarhouse/pygmy/
RUN GO111MODULE=on GOOS=linux GOARCH=amd64 go build -o pygmy-go-linux-x86 .
RUN GO111MODULE=on GOOS=darwin GOARCH=amd64 go build -o pygmy-go-darwin .

RUN ls

FROM alpine
WORKDIR /app
COPY --from=builder /go/src/github.com/fubarhouse/pygmy/pygmy-go-darwin .
COPY --from=builder /go/src/github.com/fubarhouse/pygmy/pygmy-go-linux-x86 .

