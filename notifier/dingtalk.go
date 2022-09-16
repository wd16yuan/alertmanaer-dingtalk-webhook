package notifier

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"alertmanaer-dingtalk-webhook/model"
	"alertmanaer-dingtalk-webhook/transformer"

	"github.com/coocood/freecache"
)

var (
	TokenInvalidErr = errors.New("token invalid")
	URLInvalidErr   = errors.New("dingtalk robot url invalid")
	KeyByte         = []byte("NRHp_op=K7rSI_#dft+3gQpYqlSu^VWT")
	IV              = []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f}
	cache           *freecache.Cache
	expire          = 60 * 60
)

func init() {
	cacheSize := 10 * 1024 * 1024
	cache = freecache.NewCache(cacheSize)
}

// Send send markdown message to dingtalk
func Send(notification model.Notification, defaultUrl, token, secret string) (err error) {

	markdown, robotURL, token, secret, err := transformer.TransformToMarkdown(notification)

	if err != nil {
		return
	}

	data, err := json.Marshal(markdown)
	if err != nil {
		return
	}

	var dingTalkRobotURL string

	if robotURL != "" {
		dingTalkRobotURL = robotURL
		if token == "" {
			fmt.Printf("dingtalk webhook: %s, token invalid.\n", robotURL)
			return TokenInvalidErr
		}
	} else {
		dingTalkRobotURL = defaultUrl
	}

	if len(dingTalkRobotURL) == 0 {
		return URLInvalidErr
	}

	token, err = GetCacheString(token)
	if err != nil {
		return err
	}

	if secret != "" {
		secret, err = GetCacheString(secret)
		if err != nil {
			return err
		}
		timestamp, sign := GetSignature(secret)
		dingTalkRobotURL = fmt.Sprintf("%s?access_token=%s&timestamp=%d&sign=%s", dingTalkRobotURL, token, timestamp, sign)
	} else {
		dingTalkRobotURL = fmt.Sprintf("%s?access_token=%s", dingTalkRobotURL, token)
	}

	req, err := http.NewRequest(
		"POST",
		dingTalkRobotURL,
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

func GetSignature(secret string) (timestamp int64, sign string) {
	timestamp = time.Now().UnixMilli()
	msg := fmt.Sprintf("%d\n%s", timestamp, secret)

	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(msg))
	signData := h.Sum(nil)
	sign = url.QueryEscape(base64.StdEncoding.EncodeToString(signData))
	return
}

func DecryptString(plainText string) (cipherText string, err error) {
	plainByte, err := base64.StdEncoding.DecodeString(plainText)
	if err != nil {
		return
	}

	c, err := aes.NewCipher(KeyByte)
	if err != nil {
		return
	}

	cfb := cipher.NewCFBEncrypter(c, IV)
	cipherByte := make([]byte, len(plainByte))
	cfb.XORKeyStream(cipherByte, plainByte)
	cipherText = string(cipherByte)
	return
}

func EncryptString(cipherText string) (plainText string, err error) {
	cipherByte := []byte(cipherText)
	c, err := aes.NewCipher(KeyByte)
	if err != nil {
		return
	}

	cfbdec := cipher.NewCFBDecrypter(c, IV)
	plainByte := make([]byte, len(cipherByte))
	cfbdec.XORKeyStream(plainByte, cipherByte)
	plainText = base64.StdEncoding.EncodeToString(plainByte)
	return
}

func PrintEncryptTokenAndSecret(token, secret string) {
	token_, err := EncryptString(token)
	if err != nil {
		fmt.Printf("token '%s' encrypt fail. Error: %s \n", token, err.Error())
	}
	secret_, err := EncryptString(secret)
	if err != nil {
		fmt.Printf("secret '%s' encrypt fail. Error: %s \n", secret, err.Error())
	}
	fmt.Printf("Encrypt Token: %s\nEncrypt Secret: %s\n", token_, secret_)
}

func CheckTokenAndSecret(token, secret string) (err error) {
	token_, err := DecryptString(token)
	if err != nil {
		return
	}
	cache.Set([]byte(token), []byte(token_), expire)
	secret_, err := DecryptString(secret)
	cache.Set([]byte(secret), []byte(secret_), expire)
	return err
}

func GetCacheString(key string) (value string, err error) {
	valueByte, err := cache.Get([]byte(key))
	if err != nil && err != freecache.ErrNotFound {
		return
	}
	if err == freecache.ErrNotFound {
		value, err = DecryptString(key)
		if err != nil {
			return
		}
		cache.Set([]byte(key), []byte(value), expire)
		return
	}
	value = string(valueByte)
	return
}
