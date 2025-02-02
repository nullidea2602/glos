package luaexec

import lua "github.com/yuin/gopher-lua"

func SafePreload(L *lua.LState) {
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
