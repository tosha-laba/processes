package uitools

import (
	"fmt"

	"github.com/nsf/termbox-go"
)

// Print печатает текст в заданной позиции экрана
func Print(x, y int, fColor, bColor termbox.Attribute, caption string) {
	symbols := []rune(caption)

	for i := 0; i < len(symbols); i++ {
		termbox.SetCell(x+i, y, symbols[i], fColor, bColor)
	}
}

// Printf печатает отформатированный текст в заданной позиции экрана
func Printf(x, y int, fColor, bColor termbox.Attribute, caption string, a ...interface{}) {
	Print(x, y, fColor, bColor, fmt.Sprintf(caption, a...))
}
