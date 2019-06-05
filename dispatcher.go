package main

import (
	"fmt"
	"math/rand"
	"operating-systems/processes/uitools"
	"sync"

	"github.com/nsf/termbox-go"
)

const (
	tableUpperBorder = "┌─────┬────────────┬───────────┬─────────────┬─────────┬───────────────┬───────┐"
	tableHeader      = "│PID  │Имя         │Память     │Состояние    │Время CPU│Осталось тактов│Квант  │"
	tableSeparator   = "├─────┼────────────┼───────────┼─────────────┼─────────┼───────────────┼───────┤"
	tableContent     = "│%5d│%12s│%11d│%13s│%9d│%15d│%7d│"
	tableLowerBorder = "└─────┴────────────┴───────────┴─────────────┴─────────┴───────────────┴───────┘"
)

// ProcessTable - представление таблицы процессов
type ProcessTable struct {
	table                  []Process
	first                  int
	processCounter         int
	roundRobinProcessIndex int
	currentProcessIndex    int
	currentProcess         *Process
}

// Add добавляет процесс в таблицу и увеличивает счетчик процессов
func (pt *ProcessTable) Add(process Process) {
	pt.table = append(pt.table, process)
	pt.processCounter++
}

// AddProcess добавляет сформированный процесс в таблицу
func (pt *ProcessTable) AddProcess() {
	proc := Process{Name: fmt.Sprintf("proc%d", pt.processCounter),
		Memory:        rand.Intn(MaxRAM / 256),
		MemoryBlock:   nil,
		CyclesRemains: rand.Intn(1024),
		TimeSlot:      1,
		State:         Readiness,
		PID:           pt.processCounter,
		CPUTime:       0,
		GID:           0}

	pt.Add(proc)
}

// Draw Отображает таблицу на экране в заданной области
func (pt *ProcessTable) Draw(x, y int) {
	uitools.Print(x, y, termbox.ColorWhite, termbox.ColorBlue, tableUpperBorder)
	uitools.Print(x, y+1, termbox.ColorWhite, termbox.ColorBlue, tableHeader)
	uitools.Print(x, y+2, termbox.ColorWhite, termbox.ColorBlue, tableSeparator)

	// for i, v := range pt.table {
	for i := 0; i <= 10; i++ {
		if len(pt.table) <= i+pt.first {
			break
		}
		v := pt.table[i+pt.first]
		uitools.Printf(x, y+3+i*2, termbox.ColorWhite, termbox.ColorBlue, tableContent, v.PID, v.Name, v.Memory, v.State.Stringify(), v.CPUTime, v.CyclesRemains, v.TimeSlot)
		uitools.Print(x, y+4+i*2, termbox.ColorWhite, termbox.ColorBlue, tableSeparator)
	}

	lbPos := y + 2 + (len(pt.table)-pt.first)*2

	if len(pt.table)-pt.first > 10 {
		lbPos = y + 2 + 11*2
	}

	uitools.Print(x, lbPos, termbox.ColorWhite, termbox.ColorBlue, tableLowerBorder)
}

var once sync.Once
var tableInstance *ProcessTable

// GetProcessTable предоставляет глобальный и единственный экземпляр таблицы процессов
func GetProcessTable() *ProcessTable {
	once.Do(func() {
		tableInstance = &ProcessTable{first: 0, processCounter: 0, roundRobinProcessIndex: -1}
	})
	return tableInstance
}

// InitDispatcher инициализирует таблицу процессом init
func InitDispatcher() {
	GetProcessTable().Add(Process{Name: "init",
		Memory:        0,
		MemoryBlock:   nil,
		PID:           0,
		GID:           1,
		CyclesRemains: 0,
		State:         Readiness})
}

// ScheduleProcess выбирает процесс для выполнения по схеме Round-Robin с растущими квантами времени
func ScheduleProcess() {
	// Последовательность
	// Выбор нового процесса для исполнения
	// Если до этого в таблице не было процессов, выбирается первый
	// Если процесс не блокирован, попытка доступа к памяти
	// Если ресурс доступен, текущий процесс помечается выбранным для исполнения

	pt := GetProcessTable()
	if pt.processCounter > 0 && pt.roundRobinProcessIndex == -1 {
		pt.roundRobinProcessIndex = 0
	}

	pt.currentProcess = &pt.table[pt.roundRobinProcessIndex]

	// Идентификатор группы == 1 => особый процесс init, исполняемым не помечается
	if pt.currentProcess.GID == 1 {
		return
	}

	// Если блокируется - надо дождаться выхода из блокировки, не помечается исполняемым
	if pt.currentProcess.State == Blocking {
		return
	}

	// Если ресурс памяти уже в RAM, выбирается для исполнения
	if pt.currentProcess.MemoryBlock != nil {
		pt.currentProcessIndex = pt.roundRobinProcessIndex
		pt.currentProcess.State = Execution
		return
	}

	// Попытка резервирования памяти
	block := GetMMU().Add(pt.currentProcess.Memory)
	if block != nil {
		pt.currentProcess.MemoryBlock = block
		pt.currentProcessIndex = pt.roundRobinProcessIndex
		pt.currentProcess.State = Execution
	}
}

// PerformProcess выполняет процесс, меняет его состояние, либо выгружает из оперативной памяти, либо завершает процесс
func PerformProcess() {
	// Алгоритм
	// Если процесс доступен для исполнения, исполняется

	// Если процесс завершился, то освобождение памяти без свопа

	// Если RAM занят больше, чем на половину, попытка выгрузить на диск
	// Если диск (почти) переполнен, аварийное завершение работы программы

	// Переход к другому процессу

	pt := GetProcessTable()

	// Если текущего процесса Round-Robin нет, пропускаем такт
	if pt.currentProcess == nil {
		return
	}

	// Если процесс не выбран для исполнения, пропускаем такт
	if pt.currentProcess.State != Execution {
		// Циклическое смещение указателя Round-Robin
		pt.roundRobinProcessIndex = (pt.roundRobinProcessIndex + 1) % len(pt.table)
		return
	}

	// Пересчет блоков процесса
	// Уменьшаем количество тактов
	if pt.currentProcess.CyclesRemains-pt.currentProcess.TimeSlot <= 0 {
		// Если кванта времени хватило на завершение процесса, ставим количество тактов процесса 0
		// Добавляем к процессорному времмени разность кванта и оставшихся тактов
		pt.currentProcess.CPUTime += pt.currentProcess.TimeSlot - pt.currentProcess.CyclesRemains
		pt.currentProcess.CyclesRemains = 0
	} else {
		// Иначе, уменьшаем число тактов на квант
		// Увеличиваем время CPU на квант
		// Увеличиваем квант в 2 раза
		pt.currentProcess.CyclesRemains -= pt.currentProcess.TimeSlot
		pt.currentProcess.CPUTime += pt.currentProcess.TimeSlot
		pt.currentProcess.TimeSlot <<= 1
	}

	mmu := GetMMU()
	isProcessRemoving := false

	// Процесс завершился
	if pt.currentProcess.CyclesRemains == 0 {
		mmu.Free(pt.currentProcess.MemoryBlock, false)
		// Завершение процесса
		// Удаление процесса из таблицы
		isProcessRemoving = true
	}

	// Если RAM заполнен больше, чем на половину, попытка выгрузить на диск
	if mmu.OccupiedRAM*2 > MaxRAM {
		correct := mmu.Free(pt.currentProcess.MemoryBlock, true)
		if correct {
			// Если выгрузка произошла успешно, процесс больше не связан с блоками RAM
			pt.currentProcess.MemoryBlock = nil
		} else {
			// Иначе - аварийное завершение процесса
			// Удаление процесса из таблицы
			isProcessRemoving = true
		}
	} else {
		// Если RAM мало заполнен, перевод текущего процесса в состояние готовности
		pt.currentProcess.State = Readiness
	}

	// Удаление процесса из таблицы
	if isProcessRemoving {
		pt.table = append(pt.table[:pt.currentProcessIndex], pt.table[pt.currentProcessIndex+1:]...)
		pt.roundRobinProcessIndex--
	}

	// Циклиеское смещение указателя Round-Robin
	pt.roundRobinProcessIndex = (pt.roundRobinProcessIndex + 1) % len(pt.table)
}
