package operations

import (
	"encoding/binary"
	"fmt"

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

// Create a new table in memory
func NewTable() *types.Table {
	return &types.Table{
		NumOfRows: 0,
		Pages:     [types.TableMaxPages][]byte{},
	}
}

func FreeTable(table *types.Table) {
	for i := range table.Pages {
		table.Pages[i] = []byte{}
	}
}

// SerializeRow writes a Row struct into the raw byte memory
func SerializeRow(source *types.Row, destination []byte) {

	// Using the binary.little endian since different CPU interpret interger differently
	// We are specifically mentioning the CPU to use the little endian format
	binary.LittleEndian.PutUint32(destination[types.IDOffset:], source.ID)

	copy(destination[types.UsernameOffset:types.UsernameOffset+types.UsernameSize], source.Username)
	copy(destination[types.EmailOffSet:types.EmailOffSet+types.EmailSize], source.Email)

}

func DeserializeRow(source []byte, destination *types.Row) {
	destination.ID = binary.LittleEndian.Uint32(source[types.IDOffset : types.IDOffset+types.IDSize])

	userNameBytes := source[types.UsernameOffset : types.UsernameOffset+types.UsernameSize]
	destination.Username = string(trimNull(userNameBytes))

	emailBytes := source[types.EmailOffSet : types.EmailOffSet+types.EmailSize]
	destination.Username = string(trimNull(emailBytes))
}

func Execute_insert(statement *types.Statement, table *types.Table) string {
	// Checking if table is full
	if table.NumOfRows >= types.TableMaxRows {
		return types.EXECUTE_TABLE_FULL
	}

	row := &(statement.Row)

	SerializeRow(row, table.RowSlot(table.NumOfRows))
	table.NumOfRows += 1

	return types.EXECUTE_SUCCESS
}

func Execute_select(statment *types.Statement, table *types.Table) string {
	var row types.Row

	for i := range table.NumOfRows {
		DeserializeRow(table.RowSlot(i), &row)
		fmt.Println(row)
	}

	return types.EXECUTE_SUCCESS
}
