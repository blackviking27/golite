package main

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"

	"github.com/blackviking27/golite/operations"
	"github.com/blackviking27/golite/types"
)

func execute_meta_command(command string) string {
	switch command {
	case ".exit":
		{
			fmt.Println("Shutting down...")
			os.Exit(0)
		}
	}
	return types.META_COMMAND_NOT_RECOGNIZED
}

func prepareStatement(command string, statement *types.Statement) string {
	commandParts := strings.Fields(command)
	status := types.PREPARE_FAILURE

	switch strings.ToUpper(commandParts[0]) {
	case types.STATEMENT_INSERT:
		(*statement).Type = types.STATEMENT_INSERT
		id, err := strconv.Atoi(commandParts[1])
		if err != nil {
			return status
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

func execute_statement(statement *types.Statement, table *types.Table) {
	switch statement.Type {
	case types.STATEMENT_SELECT:
		operations.Execute_select(statement, table)
	case types.STATEMENT_INSERT:
		operations.Execute_insert(statement, table)
	}
}

func main() {
	// Main REPL
	for {
		fmt.Println("Welcome to Go Lite DB")
		// Creating the in-memory table
		table := operations.NewTable() // Pointer to a table

		for {
			fmt.Println("golite>>")
			// reading user input
			scanner := bufio.NewScanner(os.Stdin)
			scanner.Scan()

			if scanner.Err() != nil {
				slog.Error("unable to parse the input")
			}

			currCommand := scanner.Text()

			// Taking in meta commands
			// commands that start with '.' like .exit

			if currCommand[0] == '.' {
				switch execute_meta_command(currCommand) {
				case types.META_SUCCESS:
					continue
				case types.META_COMMAND_NOT_RECOGNIZED:
				case types.META_FAILURE:
					fmt.Println("Unrecognized command or failure")
					continue
				}
			}

			var statement types.Statement

			switch prepareStatement(currCommand, &statement) {
			case types.PREPARE_SUCCESS:
				fmt.Println("Parsed statement")
			case types.PREPARE_FAILURE:
				fmt.Println("Unrecognized keyword at start")
				continue
			}

			execute_statement(&statement, table)
			fmt.Println("Executed")

		}

	}
}
