FROM golang:alpine as builder
LABEL stage=builder
COPY main.go /go/src/github.com/pygmystack/pygmy/
COPY go.sum /go/src/github.com/pygmystack/pygmy/
COPY go.mod /go/src/github.com/pygmystack/pygmy/
COPY cmd/ /go/src/github.com/pygmystack/pygmy/cmd/
COPY service/ /go/src/github.com/pygmystack/pygmy/service/

WORKDIR /go/src/github.com/pygmystack/pygmy/
RUN GO111MODULE=on go mod verify
RUN GO111MODULE=on GOOS=linux GOARCH=386 go build -o pygmy-go-linux-386 .
RUN GO111MODULE=on GOOS=linux GOARCH=386 CGO_ENABLED=0 go build -a -ldflags '-extldflags "-static"' -o pygmy-go-linux-386-static .
RUN GO111MODULE=on GOOS=linux GOARCH=arm go build -o pygmy-go-linux-arm .
RUN GO111MODULE=on GOOS=linux GOARCH=arm CGO_ENABLED=0 go build -a -ldflags '-extldflags "-static"' -o pygmy-go-linux-arm-static .
RUN GO111MODULE=on GOOS=linux GOARCH=arm64 go build -o pygmy-go-linux-arm64 .
RUN GO111MODULE=on GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -a -ldflags '-extldflags "-static"' -o pygmy-go-linux-arm64-static .
RUN GO111MODULE=on GOOS=linux GOARCH=amd64 go build -o pygmy-go-linux-amd64 .
RUN GO111MODULE=on GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -a -ldflags '-extldflags "-static"' -o pygmy-go-linux-amd64-static .
RUN GO111MODULE=on GOOS=darwin GOARCH=amd64 go build -o pygmy-go-darwin-amd64 .
RUN GO111MODULE=on GOOS=darwin GOARCH=arm64 go build -o pygmy-go-darwin-arm64 .
RUN GO111MODULE=on GOOS=windows GOARCH=amd64 go build -o pygmy-go.exe .

FROM alpine
WORKDIR /app
COPY --from=builder /go/src/github.com/pygmystack/pygmy/pygmy-go-linux-386 .
COPY --from=builder /go/src/github.com/pygmystack/pygmy/pygmy-go-linux-386-static .
COPY --from=builder /go/src/github.com/pygmystack/pygmy/pygmy-go-linux-arm .
COPY --from=builder /go/src/github.com/pygmystack/pygmy/pygmy-go-linux-arm-static .
COPY --from=builder /go/src/github.com/pygmystack/pygmy/pygmy-go-linux-arm64 .
COPY --from=builder /go/src/github.com/pygmystack/pygmy/pygmy-go-linux-arm64-static .
COPY --from=builder /go/src/github.com/pygmystack/pygmy/pygmy-go-linux-amd64 .
COPY --from=builder /go/src/github.com/pygmystack/pygmy/pygmy-go-linux-amd64-static .
COPY --from=builder /go/src/github.com/pygmystack/pygmy/pygmy-go-darwin-amd64 .
COPY --from=builder /go/src/github.com/pygmystack/pygmy/pygmy-go-darwin-arm64 .
COPY --from=builder /go/src/github.com/pygmystack/pygmy/pygmy-go.exe .
