package main

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"

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

			switch operations.PrepareStatement(currCommand, &statement) {
			case types.PREPARE_SUCCESS:
				fmt.Println("Parsed statement")
			case types.PREPARE_FAILURE:
				fmt.Println("Unrecognized keyword at start")
				continue
			}

			operations.ExecuteStatement(&statement, table)
			fmt.Println("Executed")

		}

	}
}
