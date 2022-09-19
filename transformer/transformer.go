package transformer

import (
	"bytes"
	"fmt"

	"alertmanaer-dingtalk-webhook/model"
)

//@function: TransformToMarkdown
//@description: 将notification转换为钉钉markdow消息格式
//@param: notification model.Notification
//@return: markdown *model.DingTalkMarkdown, robotURL, token, secret string, err error

func TransformToMarkdown(notification model.Notification) (markdown *model.DingTalkMarkdown, robotURL, token, secret string, err error) {

	groupKey := notification.GroupKey
	status := notification.Status

	annotations := notification.CommonAnnotations
	robotURL = annotations["dingtalkUrl"]
	token = annotations["dingtalkToken"]
	secret = annotations["dingtalkSecret"]

	var buffer bytes.Buffer

	buffer.WriteString(fmt.Sprintf("### 通知组%s (当前状态：%s) \n", groupKey, status))

	buffer.WriteString(fmt.Sprintf("#### 告警项：\n"))

	for _, alert := range notification.Alerts {
		annotations := alert.Annotations
		buffer.WriteString(fmt.Sprintf("##### %s\n > %s\n", annotations["summary"], annotations["description"]))
		buffer.WriteString(fmt.Sprintf("\n> 开始时间：%s\n", alert.StartsAt.Format("15:04:05")))
	}

	markdown = &model.DingTalkMarkdown{
		MsgType: "markdown",
		Markdown: &model.Markdown{
			Title: fmt.Sprintf("通知组：%s(当前状态:%s)", groupKey, status),
			Text:  buffer.String(),
		},
		At: &model.At{
			IsAtAll: false,
		},
	}

	return
}
