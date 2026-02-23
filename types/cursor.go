package types

type Cursor struct {
	Table        *Table
	RowNumber    int
	IsEndOfTable bool
}
