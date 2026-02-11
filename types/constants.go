package types

const (
	META_SUCCESS                = "META_SUCCESS"
	META_FAILURE                = "META_FAILURE"
	META_COMMAND_NOT_RECOGNIZED = "META_COMMAND_NOT_RECOGNIZED"
	STATEMENT_INSERT            = "INSERT"
	STATEMENT_SELECT            = "SELECT"

	PREPARE_SUCCESS                 = "PREPARE_SUCCESS"
	PREPARE_FAILURE                 = "PREPARE_FAILURE"
	PREPARE_NEGATIVE_ID_FAILURE     = "PREPARE_NEGATIVE_ID_FAILURE"
	PREPARE_STRING_TOO_LONG_FAILURE = "PREPARE_STRING_TOO_LONG_FAILURE"

	EXECUTE_TABLE_FULL = "EXECUTE_TABLE_FULL"
	EXECUTE_SUCCESS    = "EXECUTE_SUCCESS"
	EXECUTE_FAILURE    = "EXECUTE_FAILURE"

	// DB column size definition
	ColumnUsernameSize = 32
	ColumnEmailSize    = 255

	IDOffset       = 0
	IDSize         = 4
	UsernameOffset = IDOffset + IDSize
	UsernameSize   = ColumnUsernameSize
	EmailOffSet    = UsernameOffset + UsernameSize
	EmailSize      = ColumnEmailSize
	RowSize        = IDSize + UsernameSize + EmailSize

	PageSize      = 4096
	TableMaxPages = 100
	RowsPerPage   = PageSize / RowSize
	TableMaxRows  = TableMaxPages * RowsPerPage
)
