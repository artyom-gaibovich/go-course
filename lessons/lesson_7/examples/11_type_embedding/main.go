package main

import "fmt"

type Engine struct {
	Power int
}

func (e Engine) Start() {
	fmt.Println("engine started, power:", e.Power)
}

type Car struct {
	Engine
	Brand string
}

func main() {
	car := Car{
		Engine: Engine{Power: 120},
		Brand:  "Lada",
	}

	car.Start()
	fmt.Println("power via embedded name:", car.Engine.Power)
	fmt.Println("power via promotion:", car.Power)
}
