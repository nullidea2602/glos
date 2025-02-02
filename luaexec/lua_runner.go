package luaexec

import (
	"bufio"
	"fmt"
	"glos/glosfs"
	"os"
	"strings"

	lua "github.com/yuin/gopher-lua"
)

func Execute(content string, args []string) error {
	L := lua.NewState()
	defer L.Close()

	SafePreload(L)

	L.SetGlobal("read_file", L.NewFunction(luaReadFile))
	L.SetGlobal("list_files", L.NewFunction(luaListFiles))
	L.SetGlobal("write_file", L.NewFunction(luaWriteFile))
	L.SetGlobal("read_multiline_input", L.NewFunction(luaReadMultilineInput))
	L.SetGlobal("delete_file", L.NewFunction(luaDeleteFile))
	L.SetGlobal("set_env", L.NewFunction(luaSetEnv))
	L.SetGlobal("get_env", L.NewFunction(luaGetEnv))
	L.SetGlobal("clear_screen", L.NewFunction(luaClearScreen))

	luaTable := L.NewTable()
	for i, arg := range args {
		L.SetTable(luaTable, lua.LNumber(i+1), lua.LString(arg))
	}
	L.SetGlobal("args", luaTable)

	return L.DoString(content)
}

// Example usage in Lua:
// content = read_file("filename")
// print(content)
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

// Example usage in Lua:
// files = list_files()
// for filename, _ in pairs(files) do
//
//	print(filename)
//
// end
func luaListFiles(L *lua.LState) int {
	luaTable := L.NewTable()
	for filename := range glosfs.MemoryFS {
		L.SetTable(luaTable, lua.LString(filename), lua.LTrue)
	}
	L.Push(luaTable)
	return 1
}

// Example usage in Lua:
// write_file("filename", "content")
func luaWriteFile(L *lua.LState) int {
	filename := L.ToString(1)
	content := L.ToString(2)
	glosfs.MemoryFS[filename] = content
	return 0
}

// Example usage in Lua:
// content = read_multiline_input()
// print(content)
func luaReadMultilineInput(L *lua.LState) int {
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
	L.Push(lua.LString(content.String()))
	return 1
}

// Example usage in Lua, including error handling:
// success, error_message = delete_file("filename")
// if not success then
//
//	print(error_message)
//
// end
func luaDeleteFile(L *lua.LState) int {
	filename := L.ToString(1)
	if _, exists := glosfs.MemoryFS[filename]; !exists {
		L.Push(lua.LNil)
		L.Push(lua.LString("File not found"))
		return 2 // Return nil and error message
	}
	delete(glosfs.MemoryFS, filename)

	// Explicitly return true for success
	L.Push(lua.LTrue)
	return 1
}

// Example usage in Lua:
// set_env("variable_name", "variable_value")
// variable_value = get_env("variable_name")
// print(variable_value)
func luaSetEnv(L *lua.LState) int {
	varName := L.ToString(1)
	varValue := L.ToString(2)
	glosfs.GlosEnv[varName] = varValue
	L.Push(lua.LTrue)
	return 1
}

// Example usage in Lua:
// variable_value = get_env("variable_name")
// print(variable_value)
func luaGetEnv(L *lua.LState) int {
	varName := L.ToString(1)
	value, exists := glosfs.GlosEnv[varName]
	if !exists {
		L.Push(lua.LNil)
		return 1
	}
	L.Push(lua.LString(value))
	return 1
}

// Example usage in Lua:
// clear_screen()
func luaClearScreen(L *lua.LState) int {
	fmt.Print("\033[2J\033[H") // ANSI code: Clear screen and move cursor to home
	return 0
}
