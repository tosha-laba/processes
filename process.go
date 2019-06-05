package main

// ProcessState описывает учитываемые состояния прроцесса
type ProcessState int

const (
	// Execution свидетельствует о выполнении процесса
	Execution ProcessState = iota
	// Readiness о готовности
	Readiness
	// Blocking о блокировке прерыванием
	Blocking
)

// Stringify переводит вариант перечисления в строку
func (ps ProcessState) Stringify() string {
	switch ps {
	case Execution:
		return "Выполнение"
	case Readiness:
		return "Готовность"
	case Blocking:
		return "Блокировка"
	}

	return ""
}

// Process представляет процесс вместе с управляющим блоком
type Process struct {
	Name          string
	Memory        int
	MemoryBlock   *MemoryBlockNode
	CyclesRemains int
	TimeSlot      int
	State         ProcessState
	PID           int
	CPUTime       int
	GID           int
}
