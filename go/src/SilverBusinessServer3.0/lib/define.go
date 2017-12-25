package lib

//User : 用户
type User struct {
	UserID   int64  `form:"userID" json:"userID"`
	UserName string `form:"userName" json:"userName"`
	UserType int    `form:"userType " json:"userType"`
	RightList string `form:"rightList" json:"rightList"`
}

//表示的是非管理员
type UserAccount struct {
	UserID   int64  `form:"userID" json:"userID"`
	UserName string `form:"userName" json:"userName"`
	RightList string `form:"rightList" json:"rightList"`
}

//Device : 设备
type Device struct {
	DeviceID      int64  `form:"deviceID" json:"deviceID"`
	DeviceVmsID   string `form:"deviceVmsID" json:"deviceVmsID"`
	DeviceName    string `form:"deviceName" json:"deviceName"`
	DeviceIP      string `form:"deviceIP" json:"deviceIP"`
	RtspURL       string `form:"rtspURL" json:"rtspURL"`
	MainStreamURL string `form:"mainStreamURL" json:"mainStreamURL"`
	SubStreamURL  string `form:"subStreamURL" json:"subStreamURL"`
}

//设备列表
type DeviceList struct {
	DeviceList []Device `form:"deviceList" json:"deviceList"`
}

//mysql已经存在的设备

// Group :  设备组
type Group struct {
	GroupID    int64    `form:"groupID" json:"groupID"`
	GroupName  string   `form:"groupName" json:"groupName"`
	DeviceList []Device `form:"deviceList" json:"deviceList"`
}

// SearchPersonResult :
type SearchPersonResult struct {
	DeviceUUID string `form:"deviceUUID" json:"deviceUUID"`
	ImageURL   string `form:"imageURL" json:"imageURL"`
	ImageTime  int64  `form:"imageTime" json:"imageTime"`
}

type FaceFeature struct {
	Gender       []int `form:"gender" json:"gender"`
	Age          []int `form:"age" json:"age"`
	HairStyle    []int `form:"hairStyle" json:"hairStyle"`
	IsSpectacled []int `form:"isSpectacled" json:"isSpectacled"`
}

type BodyFeature struct {
	UpperBodyColor []int `form:"upperBodyColor" json:"upperBodyColor"`
	LowerBodyColor []int `form:"lowerBodyColor" json:"lowerBodyColor"`
	FullBodyColor  []int `form:"fullBodyColor" json:"fullBodyColor"`
	BodilyForm     []int `form:"bodilyForm" json:"bodilyForm"`
	Height         []int `form:"height" json:"height"`
}

type FeatureDesc struct {
	FeatureType  int    `form:"type" json:"type"`
	FeatureValue string `form:"value" json:"value"`
}

type SearchPersonRule struct {
	Face    FaceFeature `form:"face" json:"face"`
	Body    BodyFeature `form:"body" json:"body"`
	Feature FeatureDesc `form:"feature" json:"feature"`
}

type SearchPersonPushHttpResult struct {
	TotalNum   int                  `form:"totalNum" json:"totalNum"`
	ErrMsg     string               `form:"errMsg" json:"errMsg"`
	ResultList []SearchPersonResult `form:"resultList" json:"resultList"`
	Err        int                  `form:"err" json:"err"`
}

type Favorite struct {
	FavoriteID   int64  `form:"favoriteID" json:"favoriteID"`
	DeviceID     int64  `form:"deviceID" json:"deviceID"`
	ImageURL     string `form:"imageURL" json:"imageURL"`
	ImageTime    int64  `form:"imageTime" json:"imageTime"`
	SearchRuleID int64  `form:"searchRuleID" json:"searchRuleID"`
}

type DeviceGroup struct {
	Id   int64  `form:"id" json:"id"`
	Name string `form:"name" json:"name"`
}

type FavoriteGroup struct {
	Id         int64         `form:"id" json:"id"`
	Name       string        `form:"name" json:"name"`
	DeviceList []DeviceGroup `form:"deviceList" json:"deviceList"`
}

type Alarm struct {
	AlarmID       int64  `form:"alarmID" json:"alarmID"`
	DeviceID      int64  `form:"deviceID" json:"deviceID"`
	AlarmTime     int64  `form:"alarmTime" json:"alarmTime"`
	ImageID       int64  `form:"imageID" json:"imageID"`
	AlarmInfo     string `form:"alarmInfo" json:"alarmInfo"`
	ProcessStatus int64  `form:"processStatus" json:"processStatus"`
	ArchiveFlag   int64  `form:"archiveFlag" json:"archiveFlag"`
	Remark        string `form:"remark" json:"remark"`
}

type AlarmImage struct {
	ImageID     int64  `form:"imageID" json:"imageID"`
	ImageWidth  int64  `form:"imageWidth" json:"imageWidth"`
	ImageHeight int64  `form:"imageHeight" json:"imageHeight"`
	ImageFormat string `form:"imageFormat" json:"imageFormat"`
	ImageData   string `form:"imageData" json:"imageData"`
}

//算法运行表
type Alg struct {
	DeviceID  int64  `form:"deviceID" json:"deviceID"`
	AlgConfig string `form:"algConfig" json:"algConfig"`
	TaskID    string `form:"taskID" json:"taskID"`
}

//所有算法分析相机
type AlgDevice struct {
	DeviceID      int64  `form:"deviceID" json:"deviceID"`
	DeviceName    string `form:"deviceName" json:"deviceName"`
	RtspURL       string `form:"rtspURL" json:"rtspURL"`
	MainStreamURL string `form:"mainStreamURL" json:"mainStreamURL"`
	SubStreamURL  string `form:"subStreamURL" json:"subStreamURL"`
	DeviceIP      string `form:"device_IP" json:"device_IP"`
	AlgStatus	  int    `form:"algStatus" json:"algStatus"` //0-停止 1-在运行
	AlgConfig     string `form:"algConfig" json:"algConfig"`
} 

//Monitor 查询任务状态http://[ip]:[port]/queryTaskStatus 结果
type TaskResult struct {
	Err    int            `form:"err" json:"err"`
	ErrMsg string         `form:"errMsg" json:"errMsg"`
	Data   TaskStatusList `form:"data" json:"data"`
}
type TaskStatusList struct {
	TaskStatusList []TStatus `form:"taskStatusList" json:"taskStatusList"`
}
type TStatus struct {
	TaskID     string `form:"taskID" json:"taskID"`
	TaskStatus int    `form:"taskStatus" json:"taskStatus"`
}

//Monitor 创建任务 http://[ip]:[port]/createTask 结果
type CreateTask struct {
	Err    int    `form:"err" json:"err"`
	ErrMsg string `form:"errMsg" json:"errMsg"`
	Data   CTask  `form:"data" json:"data"`
}
type CTask struct {
	TaskID string `form:"taskID" json:"taskID"`
}

//Monitor 启动任务 http://[ip]:[port]/startTask 结果 (停止分析任务，执行结果同用)
type StartTask struct {
	Err    int    `form:"err" json:"err"`
	ErrMsg string `form:"errMsg" json:"errMsg"`
	Data   STask  `form:"data" json:"data"`
}
type STask struct {
	TaskCapacity int `form:"taskCapacity" json:"taskCapacity"`
}

//Monitor 删除分析任务 http://[ip]:[port]/deleteTask
type DeleteTask struct {
	Err    int    `form:"err" json:"err"`
	ErrMsg string `form:"errMsg" json:"errMsg"`
}

//Sword计算框架的"任务配置信息"
type TaskPara struct {
	DeviceID  int64  `form:"deviceID"  json:"deviceID"`
	StreamURL string `form:"streamURL" json:"streamURL"`
	BeginTime int64  `form:"beginTime" json:"beginTime"`
	EndTime   int64  `form:"endTime"   json:"endTime"`
	AlgConfig string `form:"algConfig" json:"algConfig"`
}

//查询所有算法分析相机的算法运行状态
type AlgStatus struct {
	DeviceID  int64 `form:"deviceID"  json:"deviceID"`
	AlgStatus int   `form:"algStatus" json:"algStatus"`
}

//查询Sword可用的任务计算能力
type TaskCapacity struct {
	Err    int      `form:"err" json:"err"`
	ErrMsg string   `form:"errMsg" json:"errMsg"`
	Data   Capacity `form:"data" json:"data"`
}
type Capacity struct {
	TaskCapacity int `form:"taskCapacity" json:"taskCapacity"`
	TaskLimit    int `form:"taskLimit" json:"taskLimit"`
}
