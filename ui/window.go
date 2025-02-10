package ui

type Window struct {
	Width  int
	Height int
	Title  string
}

func NewWindow(width, height int, title string) *Window {
	return &Window{
		Width:  width,
		Height: height,
		Title:  title,
	}
}
