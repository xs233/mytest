package lib

import (
	"io/ioutil"
	"net/http"
	"strings"
)

// PushSearchServerHttpMessage 推送搜索http服务器信息
func PushSearchServerHttpMessage(mchost, messageContent string) (result string, err error) {

	data := messageContent                                             //获取message的内容，将其赋值给data
	url := "http://" + mchost                                          //获取url地址                                               //获取url中的错误信息
	resp, err := http.Post(url, "text/plain", strings.NewReader(data)) // POST请求
	if err != nil {                                                    // 错误的情况下返回错误信息
		return "", err
	}

	//当函数执行到最后时，这些defer语句会按照逆序执行，最后该函数返回。
	defer resp.Body.Close()                //栈方式的 延迟
	body, err := ioutil.ReadAll(resp.Body) // 读取body文本数据和错误信息
	if err != nil {                        // 错误的情况下
		return "", err
	}

	bodyString := string(body)

	return bodyString, nil
}
