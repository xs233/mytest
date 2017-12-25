package lib

//import "time"

// User :
type User struct {
	UserID   int64  `form:"userID" json:"userID"`
	UserName string `form:"userName" json:"userName"`
	UserType int    `form:"userType " json:"userType"`
	PhoneNo  string `form:"phoneNo " json:"phoneNo"`
}

// AuthRecord :
type AuthRecord struct {
	RecordBatch    string `form:"recordBatch" json:"recordBatch"`
	DeviceSupplier string `form:"deviceSupplier" json:"deviceSupplier"`
	DeviceType     string `form:"deviceType" json:"deviceType"`
	DeviceModel    string `form:"deviceModel" json:"deviceModel"`
	DeviceNumber   int    `form:"deviceNumber" json:"deviceNumber"`
	OrderNumber    string `form:"orderNumber" json:"orderNumber"`
	SupportP2P     int    `form:"supportP2P" json:"supportP2P"`
	RecordUser     string `form:"recordUser" json:"recordUser"`
	RecordTime     string `form:"recordTime" json:"recordTime"`
}

// P2PRecord :
type P2PRecord struct {
	RecordBatch string `form:"recordBatch" json:"recordBatch"`
	P2PNumber   int    `form:"p2pNumber" json:"p2pNumber"`
	OrderNumber string `form:"orderNumber" json:"orderNumber"`
	RecordUser  string `form:"recordUser" json:"recordUser"`
	RecordTime  string `form:"recordTime" json:"recordTime"`
}

// P2P :
type P2P struct {
	P2PID  string `form:"p2pID" json:"p2pID"`
	Status int    `form:"status" json:"status"`
	MAC    string `form:"mac" json:"mac"`
}
