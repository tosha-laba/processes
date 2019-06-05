package main

import (
	"container/list"
	"sync"
)

const (
	// MaxRAM - 16 мегабайт
	MaxRAM = 4194304
	// MaxDiskSpace - 256 мегабайт
	MaxDiskSpace = 67108864
)

// MemoryBlockNodeType представляет тип сегмента в связном списке блоков памяти
type MemoryBlockNodeType int

const (
	// MemHole - пустой набор блоков
	MemHole MemoryBlockNodeType = iota
	// MemProcess - набор блоков процесса
	MemProcess
)

// MemoryBlockNode - представление набора блоков
type MemoryBlockNode struct {
	NodeType MemoryBlockNodeType
	Position int
	Size     int
}

// MemoryManagementUnit - представление менеджера памяти
type MemoryManagementUnit struct {
	OccupiedRAM  int
	OccupiedDisk int
	blockList    *list.List
}

// Add пытается занести в RAM фрагмент размером size блоков и возвращает указатель на сегмент или nil
func (mmu *MemoryManagementUnit) Add(size int) *MemoryBlockNode {
	// Поиск первого пустого сегмента, который может вместить size блоков
	for e := mmu.blockList.Front(); e != nil; e = e.Next() {
		val := e.Value.(*MemoryBlockNode)
		// Выбор только пустых сегментов
		if val.NodeType == MemHole {
			// Если размер сегмента равен нужному размеру, просто меняем тип сегмента
			if val.Size == size {
				val.NodeType = MemProcess
				mmu.OccupiedRAM += size
				mmu.OccupiedDisk -= size

				return val
			}
			// Если размер сегмента больше размера, делим сегмент на два
			if val.Size > size {
				reserved := &MemoryBlockNode{NodeType: MemProcess, Position: val.Position, Size: size}
				val.Position = reserved.Position + reserved.Size
				val.Size -= size
				mmu.blockList.InsertBefore(reserved, e)
				mmu.OccupiedRAM += size
				mmu.OccupiedDisk -= size

				return e.Prev().Value.(*MemoryBlockNode)
			}
		}
	}

	// nil, если не нашли сегмент
	return nil
}

// Free пытается выгрузить из RAM указанный сегмент, записав или не записав его на диск
// Возвращает флаг успешной очистки сегмента
func (mmu *MemoryManagementUnit) Free(block *MemoryBlockNode, isSwap bool) bool {
	// Итерирование по всем блокам
	for e := mmu.blockList.Front(); e != nil; e = e.Next() {
		val := e.Value.(*MemoryBlockNode)

		// Если указанный блок найден
		if block == val {
			// Если требуется выгрузка на диск и диск переполнен, неудачное выполнение операции
			if isSwap && val.Size+mmu.OccupiedDisk > MaxDiskSpace {
				return false
			}

			// Очистка сегмента
			val.NodeType = MemHole
			// Корректировка показателей памяти и дискового пространства
			mmu.OccupiedRAM -= val.Size
			if isSwap {
				mmu.OccupiedDisk += val.Size
			}

			var prevVal *MemoryBlockNode
			var nextVal *MemoryBlockNode

			if e.Prev() != nil {
				prevVal = e.Prev().Value.(*MemoryBlockNode)
			}
			if e.Next() != nil {
				nextVal = e.Next().Value.(*MemoryBlockNode)
			}

			// Если текущий элемент - первый, и имеется следующий,
			// то пытаемся объединить со следующим
			if prevVal == nil && nextVal != nil {
				if nextVal.NodeType == MemHole {
					val.Size += nextVal.Size
					mmu.blockList.Remove(e.Next())
				}
				return true
			}

			// Если текущий элемент - последний, и имеется предыдущий,
			// пытаемся объединить с предыдущим
			if prevVal != nil && nextVal == nil {
				if prevVal.NodeType == MemHole {
					prevVal.Size += val.Size
					mmu.blockList.Remove(e)
				}
				return true
			}

			//Если текущий элемент - единственный в списке, ничего не делаем
			if prevVal == nil && nextVal == nil {
				return true
			}

			// Оптимизация списка сегментов в зависимости от 4-х возможных ситуаций
			switch true {
			// Если предыдущий сегмент - сегмент процесса, а следующий - пустой, объединить со следующим
			case prevVal.NodeType == MemProcess && nextVal.NodeType == MemHole:
				val.Size += nextVal.Size
				mmu.blockList.Remove(e.Next())
				return true
			// Если предыдущий сегмент - пустой, а следующий - сегмент процесса => объединить с предыдущим
			case prevVal.NodeType == MemHole && nextVal.NodeType == MemProcess:
				prevVal.Size += val.Size
				mmu.blockList.Remove(e)
				return true
			// Если оба соседних сегмента пустые, объединить три в один сегмент
			case prevVal.NodeType == MemHole && nextVal.NodeType == MemHole:
				prevVal.Size += val.Size + nextVal.Size
				mmu.blockList.Remove(e.Next())
				mmu.blockList.Remove(e)
				return true
			// Если соседние сегменты - сегменты процесса, ничего не делать
			case prevVal.NodeType == MemProcess && nextVal.NodeType == MemProcess:
				return true
			}
		}
	}

	return false
}

var memOnce sync.Once
var mmuInstance *MemoryManagementUnit

func GetMMU() *MemoryManagementUnit {
	memOnce.Do(func() {
		mmuInstance = &MemoryManagementUnit{blockList: list.New()}
		mmuInstance.blockList.PushFront(&MemoryBlockNode{NodeType: MemHole, Position: 0, Size: MaxRAM})
	})

	return mmuInstance
}
