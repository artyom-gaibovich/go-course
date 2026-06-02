package main

import "fmt"

func main() {
	data := [...]int{1, 2, 3}
	for i, v := range data { // copy of array
		data[i] = 10
		fmt.Println(v)
	}

	for i := range &data { // not a copy of array
		data[i] = 999
		fmt.Println(i)
	}

	for i := range data[:] {
		data[i] = 111 // not a copy of array

		fmt.Println(i)
	}
	fmt.Println(data)

}
