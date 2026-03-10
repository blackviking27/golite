package operations

import (
	"encoding/binary"
	"fmt"
	"io"
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
	fileLength := uint32(fileStat.Size())
	page := &types.Pager{
		FileDescriptor: file,
		FileLength:     fileLength,
		NumberOfPages:  uint32(fileLength / types.PageSize),
		Pages:          [types.TableMaxPages][]byte{},
	}

	if fileLength%types.PageSize != 0 {
		fmt.Println("DB file does not have whhoole pages, Corrupt file", filename)
		os.Exit(1)
	}

	return page
}

// Create a new table in memory
func DbOPen(filename string) *types.Table {

	pager := PagerOpen(filename)

	if pager.NumberOfPages == 0 {
		rootNode := GetPage(pager, 0)
		InitializeLeafNode(rootNode)
	}

	return &types.Table{
		Pager:          pager,
		RootPageNumber: 0,
	}

}

// Closing the DB
func DBClose(table *types.Table) {
	pager := table.Pager

	fmt.Printf("Flushing %d full pages to disk...\n", pager.NumberOfPages)

	// flushing the full page data to the file
	for i := 0; i < int(pager.NumberOfPages); i++ {
		if pager.Pages[i] == nil {
			continue
		}

		// flushing the data to disk
		FlushPage(pager, i)

		// The previous data at this location would be cleaned the GC
		pager.Pages[i] = nil
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
func FlushPage(pager *types.Pager, pageNum int) error {

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

func GetPage(pager *types.Pager, pageNum uint32) []byte {

	// Page not found
	if pageNum > types.TableMaxPages {
		fmt.Println("Tried to fetch page that does not exit yet")
		os.Exit(1)
	}

	// Cache miss
	if pager.Pages[pageNum] == nil {
		// Create a page of size (4096 -> PageSize we defined earlier) which is a slice of bytes
		// and it is initialised to 0
		page := make([]byte, types.PageSize)

		// Finding number of pages in current file
		numberOfPages := pager.FileLength / types.PageSize

		// If there are bytes more than multiple of 4096 then we have on more page
		// example if we have 5012 bytes then 5012 // 4095 -> 1 and 5012 % 4095 != 0, which means we have another page with pending data
		if pager.FileLength%types.PageSize != 0 {
			numberOfPages += 1
		}

		// Read from the file, if the page exists in the file
		if pageNum < numberOfPages {
			offset := pageNum * types.PageSize

			// If the page exists, then we are loading the file page data into the "page" byte slice
			_, err := pager.FileDescriptor.ReadAt(page, int64(offset))

			if err != nil && err != io.EOF {
				fmt.Println("Error while reading the DB file", err.Error())
				os.Exit(1)
			}

		}

		pager.Pages[pageNum] = page
		if pageNum >= pager.NumberOfPages {
			pager.NumberOfPages += 1
		}
	}

	return pager.Pages[pageNum]

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

	node := GetPage(table.Pager, table.RootPageNumber)

	// Checking if table is full
	if GetLeafNodeNumCells(node) >= types.LeafNodeMaxCells {
		return types.EXECUTE_TABLE_FULL
	}

	row := &(statement.Row)
	cursor := TableEndCursor(table)

	LeafNodeInsert(cursor, row.ID, row)
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

	rootPageNum := table.RootPageNumber
	cellNum := 0

	rootNode := GetPage(table.Pager, rootPageNum)
	numCells := GetLeafNodeNumCells(rootNode)

	return &types.Cursor{
		Table:        table,
		PageNum:      rootPageNum,
		CellNum:      uint32(cellNum),
		IsEndOfTable: numCells == 0,
	}
}

func TableEndCursor(table *types.Table) *types.Cursor {

	pageNum := table.RootPageNumber
	rootNode := GetPage(table.Pager, pageNum)

	numCells := GetLeafNodeNumCells(rootNode)

	return &types.Cursor{
		Table:        table,
		PageNum:      pageNum,
		CellNum:      numCells,
		IsEndOfTable: true,
	}
}

func CursorAdvance(cursor *types.Cursor) {
	// cursor.RowNumber += 1
	// if cursor.RowNumber >= int(cursor.Table.NumOfRows) {
	// 	cursor.IsEndOfTable = true
	// }

	pageNum := cursor.PageNum
	node := GetPage(cursor.Table.Pager, pageNum)

	cursor.CellNum += 1
	if cursor.CellNum >= GetLeafNodeNumCells(node) {
		cursor.IsEndOfTable = true
	}

}
func CursorValue(cursor *types.Cursor) []byte {
	pageNum := cursor.PageNum

	// Fetch the page
	page := GetPage(cursor.Table.Pager, uint32(pageNum))

	return LeafNodValue(page, cursor.CellNum)

	// rowOffSet := rowNum % types.RowsPerPage
	// byteOffSet := rowOffSet * types.RowSize

	// return page[byteOffSet : byteOffSet+types.RowSize]

}
