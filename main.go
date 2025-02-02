package main

import (
	"fmt"
	"glos/glosfs"
	"glos/luaexec"
	"glos/ui"
	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	windowWidth  = 800
	windowHeight = 600
	maxInputLen  = 256
)

func main() {
	glosfs.LoadMemoryFS() // Load memoryFS on start

	rl.InitWindow(windowWidth, windowHeight, "GLOS GUI")
	defer rl.CloseWindow()
	rl.SetTargetFPS(60)

	input := ""

	for !rl.WindowShouldClose() {
		ui.DrawUI(input)

		if ui.HandleInput(&input, maxInputLen) {
			runLuaCommand(strings.ToLower(input))
			input = ""
		}
	}

}

func runLuaCommand(input string) {
	args := strings.SplitN(input, " ", 2)
	command := args[0]
	var param string
	if len(args) > 1 {
		param = args[1]
	}

	switch command {
	case "exit":
		glosfs.SaveMemoryFS() // Save memoryFS before exit
	case "run":
		runLua(param)
	default:
		if _, exists := glosfs.MemoryFS[command]; exists {
			runLua(command + " " + param) // Execute exact match
		} else if _, exists := glosfs.MemoryFS[command+".lua"]; exists {
			runLua(command + ".lua" + " " + param) // Try with .lua extension
		} else {
			ui.Output = "Unknown command"
		}
	}
}

func runLua(input string) {
	args := strings.SplitN(input, " ", 2)
	filename := args[0]
	var luaArgs []string

	if len(args) > 1 {
		luaArgs = strings.Fields(args[1])
	}

	content, exists := glosfs.MemoryFS[filename]
	if !exists {
		ui.Output = fmt.Sprintf("File '%s' not found.\n", filename)
	}

	if err := luaexec.Execute(content, luaArgs); err != nil {
		ui.Output = fmt.Sprintln("Error executing Lua:", err)
	}
}
