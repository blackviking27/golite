package types

import (
	"fmt"
	"io"
	"os"
)

type Table struct {
	NumOfRows uint32
	Pager     *Pager
}

type Pager struct {
	FileDescriptor *os.File
	FileLength     uint32
	Pages          [TableMaxPages][]byte
}

func (this *Table) GetPage(pager *Pager, pageNum uint32) []byte {

	// Page not found
	if pageNum > TableMaxPages {
		fmt.Println("Tried to fetch page that does not exit yet")
		os.Exit(1)
	}

	// Cache miss
	if pager.Pages[pageNum] == nil {
		// Create a page of size (4096 -> PageSize we defined earlier) which is a slice of bytes
		// and it is initialised to 0
		page := make([]byte, PageSize)

		// Finding number of pages in current file
		numberOfPages := pager.FileLength / PageSize

		// If there are bytes more than multiple of 4096 then we have on more page
		// example if we have 5012 bytes then 5012 // 4095 -> 1 and 5012 % 4095 != 0, which means we have another page with pending data
		if pager.FileLength%PageSize != 0 {
			numberOfPages += 1
		}

		// Read from the file, if the page exists in the file
		if pageNum < numberOfPages {
			offset := pageNum * PageSize

			// If the page exists, then we are loading the file page data into the "page" byte slice
			_, err := pager.FileDescriptor.ReadAt(page, int64(offset))

			if err != nil && err != io.EOF {
				fmt.Println("Error while reading the DB file", err.Error())
				os.Exit(1)
			}

		}

		pager.Pages[pageNum] = page
	}

	return pager.Pages[pageNum]

}

// Table Methods
func (this *Table) RowSlot(rowNum uint32) []byte {
	pageNum := rowNum / RowsPerPage

	// Fetch the page
	page := this.GetPage(this.Pager, pageNum)

	rowOffSet := rowNum % RowsPerPage
	byteOffSet := rowOffSet * RowSize

	return page[byteOffSet : byteOffSet+RowSize]

}
