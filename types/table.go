package types

import (
	"os"
)

type Table struct {
	RootPageNumber uint32
	Pager          *Pager
}

type Pager struct {
	FileDescriptor *os.File
	FileLength     uint32
	NumberOfPages  uint32
	Pages          [TableMaxPages][]byte
}
