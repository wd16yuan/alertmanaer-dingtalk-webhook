package core

import (
	"alertmanaer-dingtalk-webhook/model"
	"alertmanaer-dingtalk-webhook/notifier"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RunServer(hook *model.DingTalkHook) {
	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Welcome to dingtalk alarm sending api!")
	})
	authorized := router.Group("/", gin.BasicAuth(gin.Accounts{
		"gY725qUW": "_D6Nc62-nr=xXw4D_=-qqUrU54c=sc3k",
	}))

	authorized.POST("/webhook", func(c *gin.Context) {
		var notification model.Notification

		err := c.BindJSON(&notification)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err = notifier.Send(notification, hook)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		}

		c.JSON(http.StatusOK, gin.H{"message": "send to dingtalk successful!"})

	})
	router.Run()
}
