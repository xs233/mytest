package errcode

// 以 1xxxxx 开头的错误表示 通用错误, 比如表单验证失败, 未登录等
// 以 2xxxxx 开头的错误表示 数据错误, 比如按照指定的表单查询数据库, 但数据没有找到等.
// 其他的开发过程中遇到, 再自行增加, 加完后记得改注释

// New :
// func New(code int, text string) error {
// 	return &E{code, text}
// }

// Err :
type Err struct {
	Code   int
	String string
}

// ErrNoError : 没有错误
var ErrNoError = Err{0, "OK"}

// ErrServer : 服务端通用错误
var ErrServer = Err{10001, "Server error"}

// ErrParams : API参数出错
var ErrParams = Err{10002, "Params error"}

// ErrForm : 表单验证失败, 请查询 API 文档
var ErrForm = Err{10003, "Form error"}

// ErrNotFound : 没有该对象
var ErrNotFound = Err{10004, "Not found"}

// ErrCookie : 获取cookie失败
var ErrCookie = Err{10005, "Get cookie failed"}

// ErrMessage : 消息错误
var ErrMessage = Err{10006, "Message error"}

// ErrFakeRequest : 伪造请求
var ErrFakeRequest = Err{10007, "Fake request"}

// ErrPermissionDenied : 没有权限
var ErrPermissionDenied = Err{10008, "Permission denied"}

// ErrQuery : 查询出错
var ErrQuery = Err{10009, "Query error"}

// ErrCreate : 创建出错
var ErrCreate = Err{10010, "Create error"}

// ErrUpdate : 更新出错
var ErrUpdate = Err{10011, "Update error"}

// ErrDelete : 删除出错
var ErrDelete = Err{10012, "Delete error"}

/////////////Account Module Error Code/////////////

// ErrUnLogin : 未登录
var ErrUnLogin = Err{10100, "You must login to continue"}

// ErrAuthenticate : 鉴权未通过
var ErrAuthenticate = Err{10101, "Username or password error"}

// ErrRegister : 注册失败
var ErrRegister = Err{10102, "Register failed"}

// ErrAccountExisting :账号已存在
var ErrAccountExisting = Err{10103, "Account is existing"}

// ErrLogin : 登录失败
var ErrLogin = Err{10104, "Login failed"}

// ErrModifyPassword : 修改密码出错
var ErrModifyPassword = Err{10105, "Modify password error"}

// ErrResetPassword : 重置密码出错
var ErrResetPassword = Err{10106, "Reset password error"}

// ErrSendShortMessage : 发送短消息出错
var ErrSendShortMessage = Err{10107, "Send short message error"}

// ErrVerifyPhoneCode : 手机验证码出错
var ErrVerifyPhoneCode = Err{10108, "Verify phone code error"}

// ErrPhoneNoFormat : 手机号码格式出错
var ErrPhoneNoFormat = Err{10109, "Phone no format error"}
