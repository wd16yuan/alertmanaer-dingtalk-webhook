FROM golang:1.16.2 as builder
WORKDIR /alertmanaer-dingtalk-webhook
COPY . .
RUN GOPROXY=https://proxy.golang.com.cn,direct CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GOARM=6 go build -a -installsuffix cgo -o app .

FROM alpine:3.16.2
RUN apk --no-cache add ca-certificates && apk --no-cache add tzdata
WORKDIR /data
COPY --from=builder /alertmanaer-dingtalk-webhook/app /data/
ENTRYPOINT ["./app"]