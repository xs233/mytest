package server

const (
	Success int = iota
	AuthorizationError
	ReadBodyError
	PushError
)

const (
	url string = "https://api.jpush.cn/v3/push"
	devKey string = "4d261278b4bfe24d874c79e9"
	devSecret string = "a985785732a6c278c2bed6c1"
)

type msgHandler struct {
	Subscriber []string `json:"subscriber"`
	MsgTitle string	`json:"msgTitle"`
	MsgBody interface{} `json:"msgBody"`
}

type responseMsg struct {
	ErrCode int	`json:"err_code"`
	ErrMsg string `json:"err_msg"`
}

type jPush struct {
	Platform string `json:"platform"`
	Audience audience `json:"audience"`
	Notification map[string]notification `json:"notification,omitempty"`
	Message string `json:"message,omitempty"`
	SmsMessage string `json:"sms_message,omitempty"`
	Options options `json:"options,omitempty"`
	Cid string `json:"cid,omitempty"`
}

type audience struct {
	RegistrationId []string `json:"registration_id"`
}

type notification struct {
	Alert string `json:"alert"`
	Extras interface{} `json:"extras"`
}

type options struct {
	TimeToLive int `json:"time_to_live"`
	ApnsProduction bool `json:"apns_production"`
}