package ocf

type (
	Error int
)

const (
	ErrSuccess                 Error = 0
	ErrGeneric                 Error = 1
	ErrBadArguments            Error = 2
	ErrNotImplemented          Error = 3
	ErrInsufficientPermissions Error = 4
	ErrNotInstalled            Error = 5
	ErrMisconfigured           Error = 6
	ErrNotRunning              Error = 7
	ErrRunningMaster           Error = 8
	ErrMasterFailed            Error = 9
)

var (
	errorMessages = []string{
		`everything is OK`,
		`generic error`,
		`bad arguments`,
		`action is not implemented`,
		`insufficient permissions`,
		`not installed`,
		`wrong configuration`,
		`is not running`,
		`running master`,
		`master failed`,
	}
)

func (e Error) Code() int     { return int(e) }
func (e Error) Error() string { return errorMessages[int(e)] }
