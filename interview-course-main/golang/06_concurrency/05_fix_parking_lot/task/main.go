// Задача: Найдите и исправьте проблемы в коде ниже
// У парковки ограниченное кол-во мест для парковки - это условие стоит иметь в виду
// Программа должна отработать корректно и завершить свою работу без зависаний
//
// Дополнение к заданию в видео
// Все авто должны приехать, припарковаться и уехать в режиме очереди
// Под "имитацией времени стоянки" имеется в виду именно время, проведенное авто на парковке

package main

import (
	"fmt"
	"sync"
	"time"
)

type ParkingLot struct {
	slots int64 // Количество мест на парковке
}

func (p *ParkingLot) Park(carID int64) {
	fmt.Printf("Тачка %d припарковывается...\n", carID)
	time.Sleep(time.Second) // Имитация времени стоянки
	fmt.Printf("Тачка %d уехала с парковки.\n", carID)
}

func main() {
	parking := &ParkingLot{slots: 3}

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
	fmt.Println("Все тачки припаркованы.")
}
