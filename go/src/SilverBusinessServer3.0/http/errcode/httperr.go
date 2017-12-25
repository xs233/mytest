package errcode

// 以 1xxxxx 开头的错误表示 通用错误, 比如表单验证失败, 未登录等
// 以 2xxxxx 开头的错误表示 数据错误, 比如按照指定的表单查询数据库, 但数据没有找到等.
// 其他的开发过程中遇到, 再自行增加, 加完后记得改注释

//Err : 在这里是定义了错误的编码和错误的提示
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
var ErrQuery = Err{10010, "Query error"}

// ErrCreate : 创建出错
var ErrCreate = Err{10011, "Create error"}

// ErrUpdate : 更新出错
var ErrUpdate = Err{10012, "Update error"}

// ErrDelete : 删除出错
var ErrDelete = Err{10013, "Delete error"}

// ErrUnLogin : 未登录
var ErrUnLogin = Err{10100, "You must login to continue"}

// ErrAuthenticate : 鉴权未通过
var ErrAuthenticate = Err{10101, "Username or password error"}

/////////////Account Module Error Code/////////////

// ErrRegister : 注册失败
var ErrRegister = Err{10102, "Register failed"}

// ErrAccountExisting :账号已存在
var ErrAccountExisting = Err{10103, "Account is existing"}

// ErrLogin : 登录失败
var ErrLogin = Err{10104, "Login failed"}

// ErrModifyPassword : 修改密码出错
var ErrModifyPassword = Err{10105, "Modify password error"}

// ErrResetPassword : 重置密码错误
var ErrResetPassword = Err{10106, "Reset password error or newPassword equal oldPassword"}

// ErrAddUser : 添加用户错误
var ErrAddUser = Err{10107, "Add user error"}

// ErrAddUser : 添加报警备注错误
var ErrAddAlarmRemark = Err{10108, "Add alarm remark error"}

// ErrPrOrAr : 处理或归档报警错误
var ErrPrOrAr = Err{10109, "ProcessStatus or archiveFlag alarm error"}

// ErrAlarmImage : 查询报警图片错误
var ErrAlarmImage = Err{10110, "Query alarm image info error"}

// ErrOpenEndALG : 算法起停错误
var ErrOpenEndALG = Err{10111, "Start or stop ALG error"}

// ErrAddDidList : 添加相机错误
var ErrAddDidList = Err{10112, "Add device id list error"}

//ErrStopALG:提示请停止运行算法
var ErrStopALG = Err{10113, "Please Stop Alg before"}

//ErrConflict:提示请停止/启动算法冲突
var ErrConflict = Err{10114, "Command Conflict! please check it"}

//ErrDeleteTask:提示请停止/启动算法冲突
var ErrDeleteTask = Err{10115, "Delete task error"}

//ErrGetALG:获取alg错误
var ErrGetALG = Err{10116, "Get Alg error"}

//ErrGetALGDevice:获取alg错误
var ErrGetALGDevice = Err{10117, "Get Alg Device error"}

//ErrJudgeCapacity 在判断 sword计算能力的时候出错
var ErrJudgeCapacity = Err{10118, "Judge Capacity error"}

//ErrScarceCapacity sword 计算能力不足
var ErrScarceCapacity = Err{10119, "Scarce Capacity error"}

//ErrJudgedRunningState 判断算法运行状态时候出错
var ErrJudgedRunningState = Err{10120, "Judged running state error"}

//ErrUpdateAlgConfig 修改算法配置出错
var ErrUpdateAlgConfig = Err{10121, "Update ALG algConfig error"}

//ErrExistAlg 删相机的时候判断alg表中是否存在
var ErrExistAlg = Err{10122, "Delete device Exist Alg error"}

//ErrResetRight 重置用户权限出错
var ErrResetRight = Err{10123, "Reset right error"}

//ErrGetAlgConfig 获取算法配置出错
var ErrGetAlgConfig = Err{10124, "Get algorithm configuration error"}