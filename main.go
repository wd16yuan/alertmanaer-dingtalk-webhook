package main

import (
	"alertmanaer-dingtalk-webhook/core"
	"alertmanaer-dingtalk-webhook/initialize"
)

func main() {
	hook := initialize.ParseCommand()
	core.RunServer(hook)
}
