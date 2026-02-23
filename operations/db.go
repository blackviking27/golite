package operations

import (
	"encoding/binary"
	"fmt"
	"os"

	"github.com/blackviking27/golite/types"
)

func trimNull(b []byte) []byte {
	for i, v := range b {
		if v == 0 {
			return b[:i]
		}
	}
	return b
}

func PagerOpen(filename string) *types.Pager {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0600)

	if err != nil {
		fmt.Println(fmt.Errorf("Unable to open DB file: %s", filename))
		os.Exit(1)
	}

	fileStat, err := file.Stat()

	if err != nil {
		fmt.Println(fmt.Errorf("Unable to open DB file: %s", filename))
		os.Exit(1)
	}

	return &types.Pager{
		FileDescriptor: file,
		FileLength:     uint32(fileStat.Size()),
		Pages:          [types.TableMaxPages][]byte{},
	}
}

// Create a new table in memory
func DbOPen(filename string) *types.Table {

	pager := PagerOpen(filename)

	numRows := pager.FileLength / types.RowSize

	return &types.Table{
		NumOfRows: numRows,
		Pager:     pager,
	}

}

// Closing the DB
func DBClose(table *types.Table) {
	pager := table.Pager

	numberOfFullPages := table.NumOfRows / types.RowsPerPage

	fmt.Printf("Flushing %d full pages to disk...\n", numberOfFullPages)

	// flushing the full page data to the file
	for i := 0; i < int(numberOfFullPages); i++ {
		if pager.Pages[i] == nil {
			continue
		}

		// flushing the data to disk
		FlushPage(pager, i, types.PageSize)

		// The previous data at this location would be cleaned the GC
		pager.Pages[i] = nil
	}

	// There could be addtional rows in the last page
	additionalRows := table.NumOfRows % types.RowsPerPage
	fmt.Printf("Flushing %d additional rows to disk...\n", additionalRows)

	if additionalRows > 0 {
		pageNum := numberOfFullPages
		if pager.Pages[pageNum] != nil {
			FlushPage(pager, int(pageNum), int(additionalRows)*types.RowSize)
			pager.Pages[pageNum] = nil
		}
	}

	err := pager.FileDescriptor.Close()

	if err != nil {
		fmt.Println("Unable to file", err)
		os.Exit(1)
	}

	for i := 0; i < types.TableMaxPages; i++ {
		pager.Pages[i] = nil
	}

}

// flushing the page data to disk
func FlushPage(pager *types.Pager, pageNum int, size int) error {

	if pager.Pages[pageNum] == nil {
		return fmt.Errorf("Error while flushing data from page %d", pageNum)
	}

	offset := pageNum * types.PageSize

	_, err := pager.FileDescriptor.WriteAt(pager.Pages[pageNum], int64(offset))

	if err != nil {
		return fmt.Errorf("Error while flushing data from page %d, %v", pageNum, err)
	}

	return nil
}

func FreeTable(table *types.Table) {
	for i := range table.Pager.Pages {
		table.Pager.Pages[i] = nil
	}
}

// SerializeRow writes a Row struct into the raw byte memory
func SerializeRow(source *types.Row, destination []byte) {

	// Using the binary.little endian since different CPU interpret interger differently
	// We are specifically mentioning the CPU to use the little endian format
	binary.LittleEndian.PutUint32(destination[types.IDOffset:], source.ID)

	// Copy the source string to the destination
	// if the length is more than the space, the string is automatically truncated
	copy(destination[types.UsernameOffset:types.UsernameOffset+types.UsernameSize], source.Username)
	copy(destination[types.EmailOffSet:types.EmailOffSet+types.EmailSize], source.Email)
}

func DeserializeRow(source []byte, destination *types.Row) {
	destination.ID = binary.LittleEndian.Uint32(source[types.IDOffset : types.IDOffset+types.IDSize])

	userNameBytes := source[types.UsernameOffset : types.UsernameOffset+types.UsernameSize]
	destination.Username = string(trimNull(userNameBytes))

	emailBytes := source[types.EmailOffSet : types.EmailOffSet+types.EmailSize]
	destination.Email = string(trimNull(emailBytes))
}

func ExecuteInsert(statement *types.Statement, table *types.Table) string {
	// Checking if table is full
	if table.NumOfRows >= types.TableMaxRows {
		return types.EXECUTE_TABLE_FULL
	}

	row := &(statement.Row)
	cursor := TableEndCursor(table)

	SerializeRow(row, CursorValue(cursor))
	table.NumOfRows += 1

	return types.EXECUTE_SUCCESS
}

func ExecuteSelect(statment *types.Statement, table *types.Table) string {
	var row types.Row

	cursor := TableStartCursor(table)

	for !cursor.IsEndOfTable {
		DeserializeRow(CursorValue(cursor), &row)
		CursorAdvance(cursor)
		fmt.Println(row)
	}

	return types.EXECUTE_SUCCESS
}

func TableStartCursor(table *types.Table) *types.Cursor {
	return &types.Cursor{
		Table:        table,
		RowNumber:    0,
		IsEndOfTable: table.NumOfRows == 0,
	}
}

func TableEndCursor(table *types.Table) *types.Cursor {
	return &types.Cursor{
		Table:        table,
		RowNumber:    int(table.NumOfRows),
		IsEndOfTable: true,
	}
}

func CursorAdvance(cursor *types.Cursor) {
	cursor.RowNumber += 1
	if cursor.RowNumber >= int(cursor.Table.NumOfRows) {
		cursor.IsEndOfTable = true
	}
}

func CursorValue(cursor *types.Cursor) []byte {
	rowNum := cursor.RowNumber
	pageNum := rowNum / types.RowsPerPage

	// Fetch the page
	page := cursor.Table.GetPage(cursor.Table.Pager, uint32(pageNum))

	rowOffSet := rowNum % types.RowsPerPage
	byteOffSet := rowOffSet * types.RowSize

	return page[byteOffSet : byteOffSet+types.RowSize]

}
