package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func startRepl(ps *pageState) {
	commandsMap := cliCommands()
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Welcome to the PokeDex! ")

	for scanner.Scan() {
		userInput := scanner.Text()
		inputParts := strings.Fields(userInput)
		if len(inputParts) == 0 {
			fmt.Print("PokeDex > ")
			continue
		}

		cmd, ok := commandsMap[inputParts[0]]
		if !ok {
			fmt.Println("Command not found:", userInput)
		}

		if len(inputParts) == 1 {
			if err := cmd.callback(ps, ""); err != nil {
				fmt.Println("Command Error:", err)
			}
		} else if err := cmd.callback(ps, inputParts[1]); err != nil {
			fmt.Println("Command Error:", err)
		}

		fmt.Print("PokeDex > ")
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
}
