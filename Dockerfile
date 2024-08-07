FROM golang:1.17 as builder

WORKDIR /workspace

COPY . .
RUN GOPROXY=https://goproxy.cn,direct go mod download

RUN GOPROXY=https://goproxy.cn,direct CGO_ENABLED=0 GOOS=linux GOARCH=amd64  go build -o uitExporter main.go

FROM alpine:3.15.4

WORKDIR /

COPY --from=builder /workspace/uitExporter .

ENTRYPOINT ["/uitExporter"]