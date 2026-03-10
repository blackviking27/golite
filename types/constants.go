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

	// NODE Type
	NOTE_INTERNAL = "NODE_INTERNAL"
	NODE_LEAF     = "NODE_LEAF"

	// Common Node header layout

	// Node type
	NodeTypeSize   uint32 = 1
	NodeTypeOffset uint32 = 1

	// Is root flag
	IsRootSize   uint32 = 1
	IsRootOffset uint32 = NodeTypeOffset + NodeTypeSize

	// Parent pointer
	ParentPointerSize   uint32 = 4
	ParentPointerOffset uint32 = IsRootOffset + IsRootSize

	CommonNodeHeaderSize = NodeTypeSize + IsRootSize + ParentPointerSize

	// Leaf node header layout
	// Cell -> key value pair
	LeafNodeNumCellsSize   uint32 = 4
	LeafNodeNumCellsOffset uint32 = CommonNodeHeaderSize
	LeafNodeHeaderSize     uint32 = CommonNodeHeaderSize + LeafNodeNumCellsSize

	// Leaf node body layout
	LeafNodeKeySize   uint32 = 4
	LeafNodeKeyOffset uint32 = 0

	LeafNodeValueSize   uint32 = RowSize
	LeafNodeValueOffset uint32 = LeafNodeKeyOffset + LeafNodeKeySize

	LeafNodeCellSize uint32 = LeafNodeKeySize + LeafNodeValueSize

	// Space left for cells after header size is defined
	LeafNodeSpaceForCells uint32 = PageSize - LeafNodeHeaderSize
	LeafNodeMaxCells      uint32 = LeafNodeSpaceForCells / LeafNodeCellSize
)
