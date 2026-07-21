package main

import "fmt"

type car struct {
	color   string
	mileage int
}

func main() {
	cars := []car{
		{
			color:   "red",
			mileage: 5_000,
		},
		{
			color:   "green",
			mileage: 10_000,
		},
		{
			color:   "blue",
			mileage: 7_000,
		},
	}

	// carPtr получает указатель на первый элемент слайса.
	carPtr := &cars[0]
	carPtr.mileage += 100 // cars[0].mileage = 5100

	// ВАЖНО: append может перераспределить память слайса, если capacity недостаточен.
	// Если capacity был 3, то append создаст новый массив и скопирует элементы.
	// В этом случае carPtr будет указывать на старый массив, а cars - на новый.
	cars = append(cars, car{color: "yellow", mileage: 15_000})

	// Если произошло перераспределение, carPtr указывает на старый массив,
	// а cars[0] находится в новом массиве. Изменение carPtr не влияет на cars[0].
	carPtr.mileage += 50

	fmt.Println(cars[0].mileage, cars[0].color)
	fmt.Println(carPtr.mileage, carPtr.color)
}

// Ответ зависит от capacity слайса:
// Если capacity >= 4: cars[0] и carPtr указывают на один объект
//   Вывод: 5150 red
//           5150 red
//
// Если capacity < 4: append перераспределяет память
//   Вывод: 5100 red
//           5150 red
//
// Объяснение:
// 1. Слайс в Go - это структура, содержащая указатель на массив, длину и capacity.
// 2. Когда capacity недостаточен, append создает новый массив и копирует элементы.
// 3. carPtr указывает на элемент в старом массиве, а cars - на новый массив.
// 4. В реальности, при создании слайса литералом capacity обычно равен длине,
//    поэтому append перераспределит память, и carPtr будет указывать на старый массив.
