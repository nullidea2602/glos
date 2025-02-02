package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	lua "github.com/yuin/gopher-lua"
)

var memoryFS = make(map[string]string) // In-memory file system
const fsFilename = "memoryfs.dat"      // Persistent storage file
const utilsDir = "utils"

func main() {
	loadMemoryFS() // Load memoryFS on start
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
			saveMemoryFS() // Save memoryFS before exit
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
			if _, exists := memoryFS[command]; exists {
				runLua(command + " " + param) // Execute exact match
			} else if _, exists := memoryFS[command+".lua"]; exists {
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
	fmt.Println("Enter text (Type ':exit' on a new line to save and exit):")

	scanner := bufio.NewScanner(os.Stdin)
	var content strings.Builder

	for scanner.Scan() {
		line := scanner.Text()
		if line == ":exit" { // Use ':exit' to stop input
			break
		}
		content.WriteString(line + "\n")
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading input:", err)
		return
	}

	memoryFS[filename] = content.String()
	fmt.Printf("File '%s' saved in memory.\n", filename)
}

func listFiles() {
	if len(memoryFS) == 0 {
		fmt.Println("No files in memory.")
		return
	}
	fmt.Println("Files in memory:")
	for filename := range memoryFS {
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

	content, exists := memoryFS[filename]
	if !exists {
		fmt.Printf("File '%s' not found.\n", filename)
		return
	}

	// Create a new Lua state
	L := lua.NewState()
	defer L.Close()

	// Remove dangerous standard libraries
	safePreload(L)

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
	content, exists := memoryFS[filename]
	if !exists {
		L.Push(lua.LNil)
		L.Push(lua.LString("File not found"))
		return 2 // Return nil and error message
	}
	L.Push(lua.LString(content))
	return 1 // Return file content
}

// Save memoryFS to a file
func saveMemoryFS() {
	file, err := os.Create(fsFilename)
	if err != nil {
		fmt.Println("Error saving memoryFS:", err)
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(memoryFS); err != nil {
		fmt.Println("Error encoding memoryFS:", err)
	}
}

func loadMemoryFS() {
	// Load existing memoryFS from file
	file, err := os.Open(fsFilename)
	if err == nil {
		defer file.Close()
		decoder := json.NewDecoder(file)
		if err := decoder.Decode(&memoryFS); err != nil {
			fmt.Println("Error decoding memoryFS:", err)
		}
	}

	// Load prewritten Lua scripts if they exist
	loadPrewrittenScripts()
}

func loadPrewrittenScripts() {
	files, err := os.ReadDir(utilsDir)
	if err != nil {
		// If directory does not exist, ignore it
		return
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".lua" {
			filename := file.Name()
			if _, exists := memoryFS[filename]; !exists { // Avoid overwriting existing files
				content, err := os.ReadFile(filepath.Join(utilsDir, filename))
				if err == nil {
					memoryFS[filename] = string(content)
					fmt.Printf("Loaded prewritten script: %s\n", filename)
				} else {
					fmt.Printf("Failed to load %s: %v\n", filename, err)
				}
			}
		}
	}
}

func safePreload(L *lua.LState) {
	// Load only safe libraries
	allowedLibs := map[string]lua.LGFunction{
		"_G":     lua.OpenBase,  // Basic Lua functions (excluding os and debug)
		"table":  lua.OpenTable, // Table manipulation
		"string": lua.OpenString,
		"math":   lua.OpenMath, // Math operations
	}

	// Open only selected libraries
	for name, lib := range allowedLibs {
		L.Push(L.NewFunction(lib))
		L.Push(lua.LString(name))
		L.Call(1, 0)
	}

	//Remove dangerous functions
	L.SetGlobal("os", lua.LNil)      // Remove os library
	L.SetGlobal("io", lua.LNil)      // Remove io library (file system access)
	L.SetGlobal("debug", lua.LNil)   // Remove debug library
	L.SetGlobal("package", lua.LNil) // Remove package manipulation
}
