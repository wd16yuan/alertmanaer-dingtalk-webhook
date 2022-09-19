package initialize

import (
	"alertmanaer-dingtalk-webhook/model"
	"flag"
	"fmt"
	"os"
	"strings"
)

var (
	h          bool
	c          bool
	token      string
	secret     string
	defaultUrl string
)

func init() {
	flag.BoolVar(&h, "h", false, "help")
	flag.BoolVar(&c, "convert", false, "token or secert convert to ciphertext")
	flag.StringVar(&defaultUrl, "defaultUrl", "https://oapi.dingtalk.com/robot/send", "global dingtalk robot webhook")
	flag.StringVar(&token, "token", "", "dingtalk robot webhook token")
	flag.StringVar(&secret, "secret", "", "dingtalk robot webhook secret")
}

func ParseCommand() *model.DingTalkHook {
	flag.Parse()
	if h {
		flag.Usage()
		os.Exit(0)
	}

	token = strings.TrimSpace(token)
	secret = strings.TrimSpace(secret)

	hook := &model.DingTalkHook{
		Url:    defaultUrl,
		Token:  token,
		Secret: secret,
	}

	if c {
		hook.ConvertTokenAndSecretToCiphertext()
		os.Exit(0)
	}
	if err := hook.CheckTokenAndSecret(); err != nil {
		fmt.Printf("token or secret inspection failed. [ %s ]", err.Error())
		os.Exit(1)
	}
	return hook
}
