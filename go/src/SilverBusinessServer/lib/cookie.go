package lib

import (
	"SilverBusinessServer/env"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

var ( //缓存关键的字段
	secretKey = env.Get("httpserver.secretcookie.secretkey").(string)
	cacheDay  = int(env.Get("httpserver.secretcookie.maxdays").(int64))
)

// SetSecretCookie :
func SetSecretCookie(writer gin.ResponseWriter, key string, value string) (err error) { //设置cookie
	cookie := http.Cookie{ //这个缓存是从http里获取的缓存数据
		Name:  key,
		Value: createSignedValue(secretKey, key, value),
		Path:  "/",
	}
	http.SetCookie(writer, &cookie)
	return
}

// GetSecretCookie :根据key的不同从cookie里面获取不同的value
func GetSecretCookie(request *http.Request, key string) (value string, err error) { //获取cookie
	singelCookie, err := request.Cookie(key)
	//fmt.Println("---->", key)
	if err != nil {
		return
	}
	svalue := singelCookie.Value

	value, err = decodeSignedValue(secretKey, key, svalue, cacheDay)
	return
}

//Get current login user ID	从cookie获取登陆账号ID
func GetCurrentUser(request *http.Request) (accountID int64, err error) {
	account, err := GetSecretCookie(request, "account")
	if err != nil {
		fmt.Println("account err")
		return -1, err
	}

	accountID, err = strconv.ParseInt(account, 10, 64) //将字符串转换成int : strconv.ParseInt(account, 10, 64) 参数1表示的是字符串形式 参数2表示的是字符串的进制 参数3表示的结果返回的bit大小
	if err != nil {
		fmt.Println("accountid err")
		return -1, err
	}

	return accountID, nil
}

// GetCurrentAccount : Get current login user account 获取真实的登陆账号，不是账号ID 不需要转换成int类型
func GetCurrentAccount(request *http.Request) (account string, err error) {
	account, err = GetSecretCookie(request, "account")
	if err != nil {
		return "", err
	}

	return account, nil
}

func formatField(s string) string {
	return strconv.Itoa(len(s)) + ":" + s // strconv.Itoa 表示的是将整数转换为十进制字符串形式（即：FormatInt(i, 10) 的简写）
	/**
	// 将整数转换为十进制字符串形式（即：FormatInt(i, 10) 的简写）
	func Itoa(i int) string

	// 将字符串转换为十进制整数，即：ParseInt(s, 10, 0) 的简写）
	func Atoi(s string) (int, error)
	*/
}

//create signature 创建（新建）签名 	SHA2565加密
func createSignature(secret string, value string) string {
	key := []byte(secret)
	mac := hmac.New(sha256.New, key) //不同的加密方式，是不是就是将sha256.New替换了？？？
	mac.Write([]byte(value))
	return hex.EncodeToString(mac.Sum(nil))
}

// 创建单一值  Base64 加密
func createSignedValue(secret string, name string, value string) string {
	timestamp := strconv.Itoa(int(time.Now().Unix())) //将整型转换成十进制字符串

	b64value := base64.URLEncoding.EncodeToString([]byte(value))
	toSign := strings.Join(
		[]string{
			"2",
			formatField("0"),
			formatField(timestamp),
			formatField(name),
			formatField(b64value),
			""},
		"|")
	return toSign + createSignature(secret, toSign)
}

func consumeField(v string) (value string) {
	sli := strings.Split(v, ":")
	length, _ := strconv.Atoi(sli[0]) //将字符串转换成10进制的整型
	value = sli[1][:length]
	return
}

//解码字符串v，说白了就是解码
func decodeFiledsValue(v string) (version int, timestamp int, key string, value string, sig string) {
	sli := strings.Split(v, "|") //切割字符串，以“|”来切割字符串v 得到的sli是一个数组
	version, _ = strconv.Atoi(consumeField(sli[1]))
	timestamp, _ = strconv.Atoi(consumeField(sli[2]))
	key = consumeField(sli[3])
	value = consumeField(sli[4])
	sig = sli[5]
	return
}

//同上
func decodeSignedValue(secret string, name string, svalue string, maxCacheDay int) (string, error) {
	if strings.Count(svalue, "|") != 5 { //统计子字符串的次数，也就是字符串中有“|”的个数
		return "", errors.New("Secure cookie format error")
	}
	_, timestamp, key, value, passedSig := decodeFiledsValue(svalue)
	signedString := svalue[:len(svalue)-len(passedSig)]
	exceptSig := createSignature(secret, signedString)
	if passedSig != exceptSig {
		return "", errors.New("consistency of signal check failed")
	}
	if key == "" || key != name {
		return "", errors.New("unknown key")
	}
	if timestamp < (int(time.Now().Unix()) - maxCacheDay*86400) {
		return "", errors.New("the signature has expired")
	}
	devaluebyte, err := base64.URLEncoding.DecodeString(value)
	devalue := string(devaluebyte)
	if err != nil {
		return "", err
	}
	return devalue, nil
}
