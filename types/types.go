package types

const (
	META_SUCCESS                = "META_SUCCESS"
	META_FAILURE                = "META_FAILURE"
	META_COMMAND_NOT_RECOGNIZED = "META_COMMAND_NOT_RECOGNIZED"
	STATEMENT_INSERT            = "INSERT"
	STATEMENT_SELECT            = "SELECT"
	PREPARE_SUCCESS             = "PREPARE_SUCCESS"
	PREPARE_FAILURE             = "PREPARE_FAILURE"
)

type Statement struct {
	Type  string
	Value string
}
