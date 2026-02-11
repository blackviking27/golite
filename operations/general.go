package operations

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/blackviking27/golite/types"
)

func PrepareStatement(command string, statement *types.Statement) string {
	commandParts := strings.Fields(command)
	status := types.PREPARE_FAILURE

	switch strings.ToUpper(commandParts[0]) {
	case types.STATEMENT_INSERT:
		(*statement).Type = types.STATEMENT_INSERT
		id, err := strconv.Atoi(commandParts[1])
		if err != nil {
			return status
		}

		if id < 0 {
			return types.PREPARE_NEGATIVE_ID_FAILURE
		}

		if len(commandParts[2]) > types.ColumnUsernameSize || len(commandParts[3]) > types.ColumnEmailSize {
			return types.PREPARE_STRING_TOO_LONG_FAILURE
		}

		(*statement).Row = types.Row{
			ID:       uint32(id),
			Username: commandParts[2],
			Email:    commandParts[3],
		}
		status = types.PREPARE_SUCCESS
	case types.STATEMENT_SELECT:
		statement.Type = types.STATEMENT_SELECT
		statement.Row = types.Row{}
		status = types.PREPARE_SUCCESS
	}
	fmt.Println(commandParts[0])
	return status
}

func ExecuteStatement(statement *types.Statement, table *types.Table) {
	switch statement.Type {
	case types.STATEMENT_SELECT:
		ExecuteSelect(statement, table)
	case types.STATEMENT_INSERT:
		ExecuteInsert(statement, table)
	}
}
