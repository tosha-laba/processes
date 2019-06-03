package uitools

import (
	"github.com/nsf/termbox-go"
)

// Button представляет кнопку, отображаемую на экране
type Button struct {
	x, y    int
	caption string
	fColor  termbox.Attribute
	bColor  termbox.Attribute
	action  func(*Button)
}

// NewButton создает экземпляр структуры Button и возвращает указатель на новый экземпляр
func NewButton(x, y int, caption string, fColor, bColor termbox.Attribute, action func(*Button)) *Button {
	return &Button{x, y, caption, fColor, bColor, action}
}

// Draw отображает кнопку в позиции x, y
func (b *Button) Draw() {
	length := len([]rune(b.caption))
	for i := 0; i <= length+1; i++ {
		termbox.SetCell(b.x+i, b.y, '▓', b.fColor, b.bColor)
		termbox.SetCell(b.x+i, b.y+2, '▓', b.fColor, b.bColor)
	}

	termbox.SetCell(b.x, b.y+1, '▓', b.fColor, b.bColor)
	termbox.SetCell(b.x+length+1, b.y+1, '▓', b.fColor, b.bColor)

	Print(b.x+1, b.y+1, b.bColor, b.fColor, b.caption)
}

// CheckClick проверяет пересечение переданных координат и отображаемого прямоугольника кнопки,
// выполняет соответствующее кнопке действие
func (b *Button) CheckClick(clickX, clickY int) {
	length := len([]rune(b.caption))
	if clickX >= b.x && clickX <= b.x+length+1 && clickY >= b.y && clickY <= b.y+2 {
		b.action(b)
	}
}

// SetPosition меняет позицию кнопки на указанную
func (b *Button) SetPosition(x, y int) {
	b.x, b.y = x, y
}
