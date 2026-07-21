// Задача: Что выведет программа и почему?

package main

import "fmt"

type car struct {
	color   string
	mileage int
}

func main() {
	cars := make([]car, 3, 3)

	cars[0] = car{
		color:   "red",
		mileage: 5_000,
	}
	cars[1] = car{
		color:   "green",
		mileage: 10_000,
	}

	cars[2] = car{
		color:   "blue",
		mileage: 7_000,
	}

	carPtr := &cars[0]
	carPtr.mileage += 100

	cars = append(cars, car{color: "yellow", mileage: 15_000})
	carPtr.mileage += 50

	cars = append(cars, car{color: "yellow", mileage: 15_000})

	fmt.Println(cars[0].mileage, cars[0].color)
	fmt.Println(carPtr.mileage, carPtr.color)
}
