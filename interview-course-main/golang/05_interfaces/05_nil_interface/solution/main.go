package main

import "fmt"

type SomeStruct struct {
	Value int
}

func CheckForNil(i interface{}) {
	if i == nil {
		fmt.Println("Это nil!")
		return
	}

	fmt.Println("Это не nil!")
}

func main() {
	var s *SomeStruct
	CheckForNil(s)
}

// Ответ:
// Это не nil!
//
// Объяснение:
// 1. var s *SomeStruct создает указатель со значением nil.
// 2. Когда мы передаем s в CheckForNil, Go создает interface{} значение.
// 3. interface{} состоит из двух частей: (type, value).
// 4. При передаче *SomeStruct(nil) создается interface{} с type=*SomeStruct, value=nil.
// 5. interface{} считается nil только если ОБЕ части nil: (nil, nil).
// 6. Здесь type=*SomeStruct (не nil), value=nil, поэтому i != nil.
//
// Чтобы проверить, является ли указатель внутри interface{} nil:
// func CheckForNil(i interface{}) {
//     if i == nil {
//         fmt.Println("Это nil!")
//         return
//     }
//     // Проверяем, является ли значение внутри интерфейса nil
//     v := reflect.ValueOf(i)
//     if v.Kind() == reflect.Ptr && v.IsNil() {
//         fmt.Println("Указатель внутри интерфейса nil!")
//         return
//     }
//     fmt.Println("Это не nil!")
// }
