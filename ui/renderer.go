package ui

import rl "github.com/gen2brain/raylib-go/raylib"

type Renderer struct {
	window *Window
}

func (r *Renderer) Draw(state *AppState) {
	rl.BeginDrawing()
	rl.ClearBackground(rl.DarkGray)
	r.drawTerminal(state.Input, state.Output)
	rl.EndDrawing()
}

func (r *Renderer) drawTerminal(input, output string) {
	rl.DrawRectangle(20, 20, 760, 560, rl.Black)
	rl.DrawText(output, 30, 30, 20, rl.Green)
	rl.DrawText(input, 30, 60, 20, rl.Black)
}
