package echo

import (
	"SilverBusinessServer/env"
	"SilverBusinessServer/lib"
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

//echo消息中心服务器 ip/port配置
var EchoIP = map[string]string{
	"ip":   env.Get("echo.host").(string),
	"port": env.Get("echo.port").(string),
}

const (
	alarmCommandCode = 1
	deleteUserCode = 10
)

/*
	URL: /publish?channel=channel1&channel=channel2
	  ---------------HTTP Body Format-------------------
	  	----------------------------------------------------------------------------------------------------------------------
		|											PROTOCOL HEADER										|   PROTOCOL BODY	 |
		----------------------------------------------------------------------------------------------------------------------
		| FLAG | LENGTH | CHECKSUM | VERSION | COMMANDCODE | ERRORCODE | TEXTDATALENGTH | BINDATALENGTH | TEXTDATA | BINDATA |
		----------------------------------------------------------------------------------------------------------------------
		|  4B  |   4B   |    4B    |    4B   |     4B      |     4B    |       4B       |      4B       |  Unknown | Unknown |
		----------------------------------------------------------------------------------------------------------------------
	  ------------------------------------------------
*/
//管理员删除用户通知所有的订阅“user_IMPOWER_DELETEUSER”频道的客户端
func PublishDelete(userID int64) error {
	url := "http://" + EchoIP["ip"] + ":" + EchoIP["port"] + "/publish?channel=user_IMPOWER_DELETEUSER"
	body := Body{UserID: userID}
	Bbody, _ := json.Marshal(body)
	len := len(Bbody)
	FLAG := []byte{'I', 'M', 'P', 'O'}
	LENGTH := IntToBytes(32)
	CHECKSUM := IntToBytes(1)
	VERSION := IntToBytes(1)
	COMMANDCODE := IntToBytes(deleteUserCode)
	ERRORCODE := IntToBytes(1)
	TEXTDATALENGTH := IntToBytes(len)
	BINDATALENGTH := IntToBytes(0)
	var HTBody []byte
	HTBody = RangeSlice(FLAG, LENGTH, CHECKSUM, VERSION, COMMANDCODE, ERRORCODE, TEXTDATALENGTH, BINDATALENGTH)
	for _, v := range Bbody {
		HTBody = append(HTBody, v)
	}
	resp, err := http.Post(url, "text/plain", strings.NewReader(string(HTBody)))
	if err != nil {
		return err
	}
	by, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(by))
	return nil
}

func PushAlarm(textDataJson *lib.AlarmValue, binData []byte) error {
	url := "http://" + EchoIP["ip"] + ":" + EchoIP["port"] + "/publish?channel=user_IMPOWER_ALARM"
	textData, _ := json.Marshal(textDataJson)
	textLen := len(textData)
	binLen := len(binData)
	FLAG := []byte{'I', 'M', 'P', 'O'}
	LENGTH := IntToBytes(32)
	CHECKSUM := IntToBytes(1)
	VERSION := IntToBytes(1)
	COMMANDCODE := IntToBytes(alarmCommandCode)
	ERRORCODE := IntToBytes(1)
	TEXTDATALENGTH := IntToBytes(textLen)
	BINDATALENGTH := IntToBytes(binLen)
	var HTBody []byte
	HTBody = RangeSlice(FLAG, LENGTH, CHECKSUM, VERSION, COMMANDCODE, ERRORCODE, TEXTDATALENGTH, BINDATALENGTH)

	for _, v := range textData {
		HTBody = append(HTBody, v)
	}

	for _, v := range binData {
		HTBody = append(HTBody, v)
	}

	resp, err := http.Post(url, "text/plain", strings.NewReader(string(HTBody)))
	if err != nil {
		fmt.Println(err)
		return err
	}
	by, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(by))

	return nil
}

//int to []byte
func IntToBytes(n int) []byte {
	tmp := int32(n)
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, tmp)
	return bytesBuffer.Bytes()
}

//字节转换成整形
func BytesToInt(b []byte) int {
	bytesBuffer := bytes.NewBuffer(b)
	var tmp int32
	binary.Read(bytesBuffer, binary.BigEndian, &tmp)
	return int(tmp)
}

//遍历不定参数
func RangeSlice(args ...[]byte) []byte {
	var body []byte
	for _, arg := range args {
		for _, v := range arg {
			body = append(body, v)
		}
	}
	return body
}

type Body struct {
	UserID int64 `form:"userID" json:"userID"`
}
