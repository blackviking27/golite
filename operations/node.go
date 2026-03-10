package operations

import (
	"encoding/binary"
	"fmt"
	"os"

	"github.com/blackviking27/golite/types"
)

// Functions to access the leaf node metadata fields

func GetLeafNodeNumCells(node []byte) uint32 {
	return binary.LittleEndian.Uint32(node[types.LeafNodeNumCellsOffset : types.LeafNodeNumCellsOffset+types.LeafNodeCellSize])
}

func SetLeafNodeNumCells(node []byte, numCells uint32) {
	binary.LittleEndian.PutUint32(node[types.LeafNodeNumCellsOffset:types.LeafNodeNumCellsOffset+types.LeafNodeNumCellsSize], numCells)
}

func LeafNodeCell(node []byte, cellNum uint32) []byte {
	offset := types.LeafNodeHeaderSize + (cellNum * types.LeafNodeCellSize)
	return node[offset:]
}

func LeafNodeKey(node []byte, cellNum uint32) []byte {
	return LeafNodeCell(node, cellNum)
}

func LeafNodValue(node []byte, cellNum uint32) []byte {
	cell := LeafNodeCell(node, cellNum)

	return cell[types.LeafNodeKeySize:]
}

func InitializeLeafNode(node []byte) {
	SetLeafNodeNumCells(node, 0)
}

func LeafNodeInsert(cursor *types.Cursor, key uint32, value *types.Row) {
	// fetching the node
	node := GetPage(cursor.Table.Pager, cursor.PageNum)

	numCells := GetLeafNodeNumCells(node)

	if numCells > types.LeafNodeMaxCells {
		fmt.Println("Need to implement splitting of nodes")
		os.Exit(1)
	}

	if cursor.CellNum < numCells {
		srcStart := types.LeafNodeHeaderSize + (cursor.CellNum * types.LeafNodeCellSize)
		srcEnd := types.LeafNodeHeaderSize + (numCells * types.LeafNodeCellSize)

		destStart := srcStart + types.LeafNodeCellSize
		destEnd := srcEnd + types.LeafNodeCellSize

		copy(node[destStart:destEnd], node[srcStart:srcEnd])
	}

	SetLeafNodeNumCells(node, numCells+1)

	// Writing the key
	keyDest := LeafNodeKey(node, cursor.CellNum)
	binary.LittleEndian.PutUint32(keyDest, key)

	// Writing the value
	valueDest := LeafNodeCell(node, cursor.CellNum)
	SerializeRow(value, valueDest)

}
