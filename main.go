package main

import (
	"bufio"
	"fmt"
	"glos/glosfs"
	"glos/luaexec"
	"os"
	"strings"

	lua "github.com/yuin/gopher-lua"
)

func main() {
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
		case "write":
			writeFile(param)
		case "ls":
			listFiles()
		case "run":
			runLua(param)
		case "help":
			fmt.Println("Commands: write <filename>, ls, run <filename> [args...], exit")
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

func writeFile(filename string) {
	if filename == "" {
		fmt.Println("Usage: write <filename>")
		return
	}

	fmt.Println("Enter text (Type ':exit' to save and exit):")

	content := readMultilineInput()
	glosfs.MemoryFS[filename] = content
	fmt.Printf("File '%s' saved in memory.\n", filename)
}

func readMultilineInput() string {
	scanner := bufio.NewScanner(os.Stdin)
	var content strings.Builder

	for scanner.Scan() {
		line := scanner.Text()
		if line == ":exit" {
			break
		}
		content.WriteString(line + "\n")
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading input:", err)
	}
	return content.String()
}

func listFiles() {
	if len(glosfs.MemoryFS) == 0 {
		fmt.Println("No files in memory.")
		return
	}
	fmt.Println("Files in memory:")
	for filename := range glosfs.MemoryFS {
		fmt.Println("-", filename)
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

	// Create a new Lua state
	L := lua.NewState()
	defer L.Close()

	// Remove dangerous standard libraries
	luaexec.SafePreload(L)

	// Register custom API functions
	L.SetGlobal("read_file", L.NewFunction(luaReadFile))
	// L.SetGlobal("shutdown", L.NewFunction(luaShutdown)) // Custom shutdown function

	// Provide script arguments
	luaTable := L.NewTable()
	for i, arg := range luaArgs {
		L.SetTable(luaTable, lua.LNumber(i+1), lua.LString(arg))
	}
	L.SetGlobal("args", luaTable)

	// Execute Lua script
	if err := L.DoString(content); err != nil {
		fmt.Println("Error executing Lua:", err)
	}
}

func luaReadFile(L *lua.LState) int {
	filename := L.ToString(1) // Get the first argument from Lua
	content, exists := glosfs.MemoryFS[filename]
	if !exists {
		L.Push(lua.LNil)
		L.Push(lua.LString("File not found"))
		return 2 // Return nil and error message
	}
	L.Push(lua.LString(content))
	return 1 // Return file content
}
