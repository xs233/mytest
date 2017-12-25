package server

import (
	"testing"
	"encoding/json"
)

type msgBody struct {
	Timestamp int64 `json:"timestamp"`
	ImageUrl string `json:"imageUrl"`
	DevId string `json:"devId"`
	Type int `json:"type"`
	Id int `json:"id"`
}

func TestServer(t *testing.T) {
	msgBody := msgBody{
		Timestamp: 1511141758513,
		ImageUrl: "29084-2017-11-20-01:35:58.jpg",
		DevId: "97e0e7054b1b441e97f8d1eec1b439fc",
		Type: 1,
		Id: 29084,
	}

	body, err := json.Marshal(msgHandler{
			Subscriber: []string{"","11","13165ffa4e3186e7857","101d8559097e79f67c9","18171adc03343fa96b8","161a3797c83d5c7f73a","191e35f7e07c7e08a04","101d8559097e055d5b6","121c83f7601df5726ab","141fe1da9e92854d294"}, 
			MsgTitle: "Just for test", 
			MsgBody: msgBody,
		},
	)

	if err != nil {
		t.Error(err)
	}
	
	if err := push("http://192.168.0.113:8090/push", "impower", "taiheming", string(body)); err != nil {
		t.Error(err)
	}
}