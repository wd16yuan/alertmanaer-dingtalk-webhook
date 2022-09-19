package notifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"alertmanaer-dingtalk-webhook/model"
	"alertmanaer-dingtalk-webhook/transformer"
)

//@function: Send
//@description: 调用钉钉机器人接口发送消息
//@param: notification model.Notification, hook *model.DingTalkHook
//@return: err error

func Send(notification model.Notification, hook *model.DingTalkHook) (err error) {

	markdown, robotURL, token, secret, err := transformer.TransformToMarkdown(notification)

	if err != nil {
		return
	}

	data, err := json.Marshal(markdown)
	if err != nil {
		return
	}

	if robotURL != "" {
		hook = &model.DingTalkHook{robotURL, token, secret}
	}

	robotURL, err = hook.GetRequestUrl()
	if err != nil {
		return
	}
	fmt.Println(robotURL)
	req, err := http.NewRequest(
		"POST",
		robotURL,
		bytes.NewBuffer(data))

	if err != nil {
		fmt.Println("dingtalk robot url not found ignore:")
		return
	}

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return
	}

	defer resp.Body.Close()
	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)

	return
}
