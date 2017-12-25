package env

const (
	// SUCCEED : Succeed
	SUCCEED int = 1
	// FAILED : Failed
	FAILED int = 0
)

// ACCOUNTTYPE : Account type
type ACCOUNTTYPE int

const (
	// AccountTypeAdmin :
	AccountTypeAdmin ACCOUNTTYPE = 1
	// AccountTypeUser :
	AccountTypeUser ACCOUNTTYPE = 2
)

const (
	// SequenceNameUserID :
	SequenceNameUserID string = "SequenceNameUserID"
)

// AESSecretIV :
const AESSecretIV string = "dvg67klq1xz45gu7"

// AESSecretKey :
const AESSecretKey string = "efds45ymqk347810"

const (
	P2P_STATUS_USED   int = 0
	P2P_STATUS_UNUSED int = 1
)

const (
	P2P_UNSUPPORT int = 0
	P2P_SUPPORT   int = 1
)
