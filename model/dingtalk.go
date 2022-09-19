package model

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"net/url"
	"time"

	"github.com/coocood/freecache"
)

var (
	KeyByte         = []byte("NRHp_op=K7rSI_#dft+3gQpYqlSu^VWT")
	IV              = []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f}
	cache           *freecache.Cache
	expire          = 60 * 60
	ErrTokenInvalid = errors.New("token invalid")
	ErrURLInvalid   = errors.New("dingtalk robot url invalid")
)

func init() {
	cacheSize := 10 * 1024 * 1024
	cache = freecache.NewCache(cacheSize)
}

type DingTalkMessage struct {
}

type At struct {
	AtMobiles []string `json:"atMobiles"`
	IsAtAll   bool     `json:"isAtAll"`
}

type DingTalkMarkdown struct {
	MsgType  string    `json:"msgtype"`
	At       *At       `json:"at"`
	Markdown *Markdown `json:"markdown"`
}

type Markdown struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}

type DingTalkHook struct {
	Url    string
	Token  string
	Secret string
}

//@function: ConvertTokenAndSecretToCiphertext
//@description: 将token、secret转换至密文并在终端显示
//@param: nil
//@return: nil

func (h *DingTalkHook) ConvertTokenAndSecretToCiphertext() {
	fmt.Println("Convert plaintext to ciphertext:")
	for _, value := range []string{h.Token, h.Secret} {
		if value == "" {
			fmt.Printf("'%s' invalid value. \n", value)
			continue
		}
		value_, err := EncryptString(value)
		if err != nil {
			fmt.Printf("\"%s\" encrypt fail. Error: \"%s\" \n", value, err.Error())
			continue
		}
		fmt.Printf("\"%s\" ======>>>  \"%s\"\n", value, value_)
	}
}

//@function: CheckTokenAndSecret
//@description: 启动时，检查token、secret是否有效，并添加至缓存
//@param: nil
//@return: err error

func (h *DingTalkHook) CheckTokenAndSecret() (err error) {
	for _, value := range []string{h.Token, h.Secret} {
		var value_ string
		value_, err = DecryptString(value)
		if err != nil {
			return
		}
		err = cache.Set([]byte(value), []byte(value_), expire)
		if err != nil {
			return
		}
	}
	return
}

//@function: GetRequestUrl
//@description: 获取钉钉机器人完整webhook
//@param: nil
//@return: robotURL string, err error

func (h *DingTalkHook) GetRequestUrl() (robotURL string, err error) {
	if h.Token == "" {
		fmt.Printf("dingtalk webhook: %s, token invalid.\n", h.Url)
		err = ErrTokenInvalid
		return
	}

	if len(h.Url) == 0 {
		err = ErrURLInvalid
		return
	}

	token, err := GetCacheString(h.Token)
	if err != nil {
		return
	}

	if h.Secret != "" {
		var secret string
		secret, err = GetCacheString(h.Secret)
		if err != nil {
			return
		}

		timestamp, sign := h.GetSignature(secret)
		robotURL = fmt.Sprintf("%s?access_token=%s&timestamp=%d&sign=%s", h.Url, token, timestamp, sign)
	} else {
		robotURL = fmt.Sprintf("%s?access_token=%s", h.Url, token)
	}
	return
}

//@function: GetSignature
//@description: 钉钉机器人secret签名
//@param: secret string
//@return: timestamp int64, sign string

func (h *DingTalkHook) GetSignature(secret string) (timestamp int64, sign string) {
	timestamp = time.Now().UnixNano() / 1e6
	msg := fmt.Sprintf("%d\n%s", timestamp, secret)

	hash := hmac.New(sha256.New, []byte(secret))
	hash.Write([]byte(msg))
	signData := hash.Sum(nil)
	sign = url.QueryEscape(base64.StdEncoding.EncodeToString(signData))
	return
}

//@function: GetCacheString
//@description: 获取缓存中的token、secret
//@param: key string
//@return: value string, err error

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

//@function: DecryptString
//@description: 解密字符串
//@param: cipherText string
//@return: plainText string, err error

func DecryptString(cipherText string) (plainText string, err error) {
	cipherByte, err := base64.StdEncoding.DecodeString(cipherText)
	if err != nil {
		return
	}

	c, err := aes.NewCipher(KeyByte)
	if err != nil {
		return
	}

	cfbdec := cipher.NewCFBDecrypter(c, IV)
	plainByte := make([]byte, len(cipherByte))
	cfbdec.XORKeyStream(plainByte, cipherByte)
	plainText = string(plainByte)
	return
}

//@function: EncryptString
//@description: 加密字符串
//@param: plainText string
//@return: cipherText string, err error

func EncryptString(plainText string) (cipherText string, err error) {
	plainByte := []byte(plainText)
	c, err := aes.NewCipher(KeyByte)
	if err != nil {
		return
	}

	cfb := cipher.NewCFBEncrypter(c, IV)
	cipherByte := make([]byte, len(plainByte))
	cfb.XORKeyStream(cipherByte, plainByte)
	cipherText = base64.StdEncoding.EncodeToString(cipherByte)
	return
}
