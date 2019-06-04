package main

import (
	"log"

	"operating-systems/processes/uitools"

	"github.com/nsf/termbox-go"
)

// TODO:
// Отрисовка таблицы процессов

func main() {
	// ==== Инициализация ресурсов библиотеки псевдографики ==== //
	if err := initTermbox(); err != nil {
		log.Fatal(err)
	}
	defer termbox.Close()

	// =============== Рабочие переменные ============== //
	// Флаг завершения работы - нужен для завершения программного цикла
	isQuitEvent := false

	// ========== Инициализация элементов управления =========== //
	createProcessButton := uitools.NewButton(1, 1, "Создать процесс", termbox.ColorWhite, termbox.ColorBlue,
		func(b *uitools.Button) {
			b.SetPosition(10, 20)
			isQuitEvent = true
		})

	for {
		// ======== Выполнение действий по нажатию ЛКМ, Escape ======== //
		pollEvents(
			func(ev *termbox.Event) {
				createProcessButton.CheckClick(ev.MouseX, ev.MouseY)
			},
			func(ev *termbox.Event) {
				isQuitEvent = true
			})
		if isQuitEvent {
			break
		}

		// ====================== Логика модели ====================== //

		// ================= Отрисовка псевдографики ================= //
		drawGUI(func() {
			createProcessButton.Draw()
		},
			termbox.ColorBlue,
			termbox.ColorBlue)
	}
}
