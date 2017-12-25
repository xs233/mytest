package env

// 在这里定义一些常量

const (
	SUCCESS int = 1
	FAILED  int = 0
)

type ACCOUNTTYPE int //这是定义一个账号的类型  1-管理员		2-用户

const (
	AccountTypeAdmin ACCOUNTTYPE = 1 //管理员
	AccountTypeUser  ACCOUNTTYPE = 2 //用户
)

const (
	// SequenceNameUserID :
	SequenceNameUserID string = "SequenceNameUserID"
	// SequenceNameGroupID :
	SequenceNameGroupID string = "SequenceNameGroupID"
	// SequenceNameSearchRuleID :
	SequenceNameSearchRuleID string = "SequenceNameRuleID"
	// SequenceNameFavoriteID :
	SequenceNameFavoriteID string = "SequenceNameFavoriteID"
)
