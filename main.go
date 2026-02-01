package main

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"
	"strings"

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
	commandType := strings.Fields(command)[0]
	status := types.PREPARE_FAILURE

	switch strings.ToUpper(commandType) {
	case types.STATEMENT_INSERT:
		statement.Type = types.STATEMENT_INSERT
		statement.Value = command
		status = types.PREPARE_SUCCESS
	case types.STATEMENT_SELECT:
		statement.Type = types.STATEMENT_SELECT
		statement.Value = command
		status = types.PREPARE_SUCCESS
	}
	fmt.Println(commandType)
	return status
}

func execute_statement(statement *types.Statement) {
	switch statement.Type {
	case types.STATEMENT_SELECT:
		fmt.Println("Select statement executed")
	case types.STATEMENT_INSERT:
		fmt.Println("Insert statement will executed")
	}
}

func main() {
	// Main REPL
	for {
		fmt.Println("Go Lite")

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

			execute_statement(&statement)
			fmt.Println("Executed")

		}

	}
}
