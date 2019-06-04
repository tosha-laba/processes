package main

import "github.com/nsf/termbox-go"

// Инициализирует библиотеку псевдографики, если инициализация прошла успешно,
// включает режим ввода с мыши и клавиши Escape.
func initTermbox() error {
	if err := termbox.Init(); err != nil {
		return err
	}

	termbox.SetInputMode(termbox.InputEsc | termbox.InputMouse)
	return nil
}

// Опрашивает события библиотеки псевдографики и выполняет переданные замыкания
func pollEvents(mouseLeftAction func(*termbox.Event), KeyEscAction func(*termbox.Event), keyboardAction func(*termbox.Event)) {
	switch ev := termbox.PollEvent(); ev.Type {
	case termbox.EventMouse:
		if ev.Key == termbox.MouseLeft {
			mouseLeftAction(&ev)
		}
	case termbox.EventKey:
		if ev.Key == termbox.KeyEsc {
			KeyEscAction(&ev)
		} else {
			keyboardAction(&ev)
		}
	}
}

// Очищает экран заданным цветом, выполняет конвеер, переданный в замыкании и выводит на экран
func drawGUI(pipeline func(), fColor, bColor termbox.Attribute) {
	termbox.Clear(fColor, bColor)
	pipeline()
	termbox.Flush()
}
