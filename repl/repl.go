package repl

import (
	"bufio"
	"fmt"
	"glos/glosfs"
	"glos/luaexec"
	"os"
	"strings"
)

// StartREPL starts the interactive shell for GLOS.
func StartREPL() {
	glosfs.LoadMemoryFS() // Load memoryFS on start
	fmt.Println("Welcome to GLOS (Go/Lua OS) 1.0")
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("glos> ")
		if !scanner.Scan() {
			break
		}
		input := scanner.Text()
		args := strings.SplitN(input, " ", 2)
		command := args[0]
		var param string
		if len(args) > 1 {
			param = args[1]
		}

		switch command {
		case "exit":
			fmt.Println("Exiting GLOS...")
			glosfs.SaveMemoryFS() // Save memoryFS before exit
			return
		case "run":
			runLua(param)
		default:
			if _, exists := glosfs.MemoryFS[command]; exists {
				runLua(command + " " + param) // Execute exact match
			} else if _, exists := glosfs.MemoryFS[command+".lua"]; exists {
				runLua(command + ".lua" + " " + param) // Try with .lua extension
			} else {
				fmt.Println("Unknown command")
			}
		}
	}
}

func runLua(input string) {
	if input == "" {
		fmt.Println("Usage: run <filename> [args...]")
		return
	}

	args := strings.SplitN(input, " ", 2)
	filename := args[0]
	var luaArgs []string

	if len(args) > 1 {
		luaArgs = strings.Fields(args[1])
	}

	content, exists := glosfs.MemoryFS[filename]
	if !exists {
		fmt.Printf("File '%s' not found.\n", filename)
		return
	}

	if err := luaexec.Execute(content, luaArgs); err != nil {
		fmt.Println("Error executing Lua:", err)
	}
}
