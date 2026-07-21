package main

import (
	"fmt"
	"sync"
	"time"
)

type ParkingLot struct {
	slots    int64
	mu       sync.Mutex
	occupied int64 // Количество занятых мест
}

func NewParkingLot(slots int64) *ParkingLot {
	return &ParkingLot{
		slots:    slots,
		occupied: 0,
	}
}

// Park паркует автомобиль, ожидая свободное место.
func (p *ParkingLot) Park(carID int64) {
	// Ждем свободное место
	for {
		p.mu.Lock()
		if p.occupied < p.slots {
			// Есть свободное место, занимаем его
			p.occupied++
			p.mu.Unlock()
			break
		}
		p.mu.Unlock()
		// Мест нет, ждем немного и проверяем снова
		time.Sleep(100 * time.Millisecond)
	}

	fmt.Printf("Тачка %d припарковалась (занято мест: %d/%d)\n", carID, p.occupied, p.slots)

	// Имитация времени стоянки
	time.Sleep(time.Second)

	// Освобождаем место
	p.mu.Lock()
	p.occupied--
	p.mu.Unlock()

	fmt.Printf("Тачка %d уехала с парковки (свободно мест: %d/%d)\n", carID, p.slots-p.occupied, p.slots)
}

func main() {
	parking := NewParkingLot(3)

	var wg sync.WaitGroup

	carIDs := []int64{1, 2, 3, 4, 5, 6}

	for _, carID := range carIDs {
		wg.Add(1)
		go func(id int64) {
			defer wg.Done()
			parking.Park(id)
		}(carID)
	}

	wg.Wait()
	fmt.Println("Все тачки припаркованы и уехали.")
}

// Суть реализации:
// 1. Используем счетчик occupied для отслеживания количества занятых мест.
// 2. Мьютекс защищает счетчик от race condition при конкурентном доступе.
// 3. Реализуем цикл ожидания свободного места: горутина ждет, пока не освободится место.
// 4. При парковке увеличиваем occupied, при отъезде уменьшаем.
// 5. Это гарантирует, что одновременно паркуется не более slots машин.
//
// Альтернативные подходы:
// - Использовать канал с буфером размером slots (semaphore pattern)
// - Использовать sync.Cond для более эффективного ожидания
