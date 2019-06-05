package main

import (
	"log"
	"math/rand"
	"time"

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
	// Индекс выделенного процесса
	selectedProcessIndex := -1
	// Координаты мыши
	hoverX, hoverY := 0, 0

	// ========== Инициализация состояния модели =========== //
	// Экземпляр таблицы процессов
	processTable := GetProcessTable()
	// Экземпляр менеджера памяти
	memoryManagementUnit := GetMMU()
	InitDispatcher()

	rand.Seed(time.Now().Unix())

	// ========== Инициализация элементов управления =========== //
	createProcessButton := uitools.NewButton(1, 1, "Создать процесс", termbox.ColorWhite, termbox.ColorBlue,
		func(b *uitools.Button) {
			processTable.AddProcess()
		})

	blockProcessButton := uitools.NewButton(20, 1, "Блокировать процесс    ", termbox.ColorWhite, termbox.ColorBlue,
		func(b *uitools.Button) {
			if selectedProcessIndex != -1 && len(processTable.table) > selectedProcessIndex {
				// Проверка, что процесс - пользовательский
				if processTable.table[selectedProcessIndex].GID != 1 {
					processTable.table[selectedProcessIndex].State = Blocking
				}
			}
		})

	unblockProcessButton := uitools.NewButton(47, 1, "Разблокировать процесс    ", termbox.ColorWhite, termbox.ColorBlue,
		func(b *uitools.Button) {
			if selectedProcessIndex != -1 && len(processTable.table) > selectedProcessIndex {
				if processTable.table[selectedProcessIndex].State == Blocking {
					processTable.table[selectedProcessIndex].State = Readiness
					// Если процесс блокировался, то он добавляется в очередь следующим для исполнения
				}
			}
		})

	for {
		// ======== Выполнение действий по нажатию ЛКМ, Escape ======== //
		pollEvents(
			func(ev *termbox.Event) {
				switch uiState {
				case ProcessMonitor:
					createProcessButton.CheckClick(ev.MouseX, ev.MouseY)
					blockProcessButton.CheckClick(ev.MouseX, ev.MouseY)
					unblockProcessButton.CheckClick(ev.MouseX, ev.MouseY)

					// Выбор элемента из таблицы процессов
					if ev.MouseX >= 0 && ev.MouseX <= 63 && ev.MouseY >= 8 && ev.MouseY <= 28 {
						if ev.MouseY%2 == 0 && len(processTable.table) > (ev.MouseY-8)/2 {
							selectedProcessIndex = (ev.MouseY-8)/2 + processTable.first
						}
					} else {
						selectedProcessIndex = -1
					}
				case MemoryDispatchMonitor:

				}

			},
			func(ev *termbox.Event) {
				if uiState == ProcessMonitor {
					switch ev.Key {
					case termbox.MouseWheelUp:
						if processTable.first > 0 {
							processTable.first--
						}
					case termbox.MouseWheelDown:
						if len(processTable.table)-processTable.first > 11 {
							processTable.first++
						}
					}
				} else {
					hoverX, hoverY = ev.MouseX, ev.MouseY
				}
			},
			func(ev *termbox.Event) {
				isQuitEvent = true
			},
			func(ev *termbox.Event) {
				switch ev.Key {
				case termbox.KeyF2:
					uiState = ProcessMonitor
				case termbox.KeyF4:
					uiState = MemoryDispatchMonitor
				case termbox.KeyArrowUp:
					if processTable.first > 0 && uiState == ProcessMonitor {
						processTable.first--
					}
				case termbox.KeyArrowDown:
					if len(processTable.table)-processTable.first > 11 && uiState == ProcessMonitor {
						processTable.first++
					}
				}
			})
		if isQuitEvent {
			break
		}

		// ====================== Логика модели ====================== //
		ScheduleProcess()

		// ================= Отрисовка псевдографики ================= //
		drawGUI(func() {
			switch uiState {
			case ProcessMonitor:
				createProcessButton.Draw()
				blockProcessButton.Draw()
				unblockProcessButton.Draw()

				if selectedProcessIndex != -1 {
					uitools.Printf(40, 2, termbox.ColorBlue, termbox.ColorWhite, "%3d", selectedProcessIndex)
					uitools.Printf(70, 2, termbox.ColorBlue, termbox.ColorWhite, "%3d", selectedProcessIndex)
				}

				processTable.Draw(0, 4)

				termbox.SetCell(80, processTable.roundRobinProcessIndex*2+7, '<', termbox.ColorBlue, termbox.ColorWhite)

				drawStatusBar()

			case MemoryDispatchMonitor:
				blockType := MemHole
				blockPos := -1
				blockSize := -1

				uitools.Print(0, 0, termbox.ColorWhite, termbox.ColorBlue, "Менеджер ресурсов")

				for i, e := 0, memoryManagementUnit.blockList.Front(); e != nil; i, e = i+1, e.Next() {
					v := e.Value.(*MemoryBlockNode)

					var blockSymbol rune
					if v.NodeType == MemProcess {
						blockSymbol = '▓'
					} else {
						blockSymbol = '░'
					}

					termbox.SetCell(1+(i%40)*2, 2+(i/40)*2, blockSymbol, termbox.ColorWhite, termbox.ColorBlue)
					if hoverX == 2+(i%40)*2 && hoverY == 3+(i/40)*2 {
						blockType = v.NodeType
						blockPos = v.Position
						blockSize = v.Size
					}

					if i != memoryManagementUnit.blockList.Len()-1 {
						termbox.SetCell(2+(i%40)*2, 2+(i/40)*2, '→', termbox.ColorWhite, termbox.ColorBlue)
					}
				}

				drawStatusBar()
				if blockPos != -1 {
					var btString string
					if blockType == MemHole {
						btString = "Пустой сегмент"
					} else {
						btString = "Сегмент процесса"
					}
					uitools.Printf(statusBarX+1, statusBarY, termbox.ColorBlue, termbox.ColorWhite, "%s Начало: %d Размер: %d", btString, blockPos, blockSize)
				}
			}
		},
			termbox.ColorBlue,
			termbox.ColorBlue)

		// ====================== Логика модели ====================== //
		PerformProcess()
	}
}
