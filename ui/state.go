package ui

type AppState struct {
	Input  string
	Output string
}

func NewAppState() *AppState {
	return &AppState{
		Output: "Welcome to GLOS GUI!\n> ",
	}
}
