package luaexec

import (
	lua "github.com/yuin/gopher-lua"
)

func Execute(content string, args []string) error {
	L := lua.NewState()
	defer L.Close()

	SafePreload(L)

	luaTable := L.NewTable()
	for i, arg := range args {
		L.SetTable(luaTable, lua.LNumber(i+1), lua.LString(arg))
	}
	L.SetGlobal("args", luaTable)

	return L.DoString(content)
}
