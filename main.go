package main

import (
	"log"

	"operating-systems/processes/uitools"

	"github.com/nsf/termbox-go"
)

// TODO:
// Отрисовка таблицы процессов

// UIStateEnum перечисляет все состояния графического интерфейса
type UIStateEnum int

const (
	// ProcessMonitor указывает, что отображается диспетчер задач
	ProcessMonitor UIStateEnum = iota
	// MemoryDispatchMonitor указывает, что отображается менеджер памяти
	MemoryDispatchMonitor
)

func main() {
	// ==== Инициализация ресурсов библиотеки псевдографики ==== //
	if err := initTermbox(); err != nil {
		log.Fatal(err)
	}
	defer termbox.Close()

	// =============== Рабочие переменные ============== //
	// Флаг завершения работы - нужен для завершения программного цикла
	isQuitEvent := false
	// Текущий экран графического интерфейса
	uiState := ProcessMonitor

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
			},
			func(ev *termbox.Event) {
				switch ev.Key {
				case termbox.KeyF1:
					uiState = ProcessMonitor
				case termbox.KeyF2:
					uiState = MemoryDispatchMonitor
				}
			})
		if isQuitEvent {
			break
		}

		// ====================== Логика модели ====================== //

		// ================= Отрисовка псевдографики ================= //
		drawGUI(func() {
			switch uiState {
			case ProcessMonitor:
				createProcessButton.Draw()

				drawStatusBar()

			case MemoryDispatchMonitor:

			}
		},
			termbox.ColorBlue,
			termbox.ColorBlue)
	}
}
