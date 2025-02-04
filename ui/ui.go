package ui

import rl "github.com/gen2brain/raylib-go/raylib"

var Output = "Welcome to GLOS GUI!\n> "

func DrawUI(input string) {
	rl.BeginDrawing()
	rl.ClearBackground(rl.DarkGray)

	//rl.DrawText("GLOS GUI Terminal", 20, 20, 20, rl.Black)
	//rl.DrawRectangleLines(20, 50, 760, 40, rl.Gray)
	//rl.DrawText(input, 30, 60, 20, rl.Black)
	rl.DrawRectangle(20, 20, 760, 560, rl.Black)
	rl.DrawText(Output, 30, 30, 20, rl.Green)

	rl.EndDrawing()
}

func HandleInput(input *string, maxInputLen int) bool {
	if rl.IsKeyPressed(rl.KeyEnter) {
		Output += "\n" // + *input
		return true    // Indicates that Enter was pressed
	}

	if rl.IsKeyPressed(rl.KeyBackspace) && len(*input) > 0 {
		*input = (*input)[:len(*input)-1]
		Output = Output[:len(Output)-1]
	}

	for char := rl.GetCharPressed(); char > 0; char = rl.GetCharPressed() {
		if len(*input) < maxInputLen {
			*input += string(char)
			Output += string(char)
		}
	}
	return false
}
