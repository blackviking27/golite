package operations

import (
	"encoding/binary"

	"github.com/blackviking27/golite/types"
)

// Create a new table in memory
func NewTable() *types.Table {
	return &types.Table{
		NumOfRows: 0,
	}
}

func trimNull(b []byte) []byte {
	for i, v := range b {
		if v == 0 {
			return b[:i]
		}
	}
	return b
}

// SerializeRow writes a Row struct into the raw byte memory
func SerializeRow(source types.Row, destination []byte) {

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
