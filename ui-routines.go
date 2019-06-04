package main

import "github.com/nsf/termbox-go"

const (
	statusBarWidth = 120
	statusBarX     = 0
	statusBarY     = 29
)

func drawStatusBar() {
	for i := 0; i < statusBarWidth; i++ {
		termbox.SetCell(statusBarX+i, statusBarY, '▓', termbox.ColorWhite, termbox.ColorWhite)
	}
}
