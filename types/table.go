package types

type Table struct {
	NumOfRows uint32
	Pages     [TableMaxPages][]byte
}

// Table Methods
func (this *Table) RowSlot(rowNum uint32) []byte {
	pageNum := rowNum / RowsPerPage

	// Check if a page does exists
	if this.Pages[pageNum] == nil {
		this.Pages[pageNum] = make([]byte, PageSize)
	}

	rowOffSet := rowNum % RowsPerPage
	byteOffSet := rowOffSet * RowSize

	return this.Pages[pageNum][byteOffSet : byteOffSet+RowSize]
}
