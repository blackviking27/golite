package types

type Cursor struct {
	Table        *Table
	PageNum      uint32
	CellNum      uint32
	IsEndOfTable bool
}
