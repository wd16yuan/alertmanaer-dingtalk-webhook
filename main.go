package main

import (
	"flag"
	"fmt"
	"net/http"

	model "alertmanaer-dingtalk-webhook/model"
	"alertmanaer-dingtalk-webhook/notifier"

	"github.com/gin-gonic/gin"
)

var (
	h          bool
	p          bool
	token      string
	secret     string
	defaultUrl string
)

func init() {
	flag.BoolVar(&h, "h", false, "help")
	flag.BoolVar(&p, "print", false, "print encrypted token and secert")
	flag.StringVar(&defaultUrl, "defaultUrl", "https://oapi.dingtalk.com/robot/send", "global dingtalk robot webhook")
	flag.StringVar(&token, "token", "", "dingtalk robot webhook token")
	flag.StringVar(&secret, "secret", "", "dingtalk robot webhook secret")
}

func main() {

	flag.Parse()
	if h {
		flag.Usage()
		return
	}
	if p {
		notifier.PrintEncryptTokenAndSecret(token, secret)
		return
	}
	if err := notifier.CheckTokenAndSecret(token, secret); err != nil {
		fmt.Printf("token or secret format error. [ %s ]", err.Error())
		return
	}

	router := gin.Default()
	router.POST("/webhook", func(c *gin.Context) {
		var notification model.Notification

		err := c.BindJSON(&notification)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err = notifier.Send(notification, defaultUrl, token, secret)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		}

		c.JSON(http.StatusOK, gin.H{"message": "send to dingtalk successful!"})

	})
	router.Run()
}
