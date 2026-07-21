package main

import (
	"fmt"
)

// appendLenWrong — НЕПРАВИЛЬНАЯ реализация.
// Функция получает копию слайса и изменяет только локальную копию.
// Вызывающая функция не увидит изменений.
func appendLenWrong(numbers []*int) {
	size := len(numbers) // size = 3
	// append изменяет локальную переменную numbers.
	// В main слайс numbers остается без изменений.
	numbers = append(numbers, &size)
}

// appendLen — правильная реализация через возврат слайса.
// Функция возвращает новый слайс, вызывающая функция сохраняет результат.
// Это идиоматичный способ в Go для функций, изменяющих слайсы.
func appendLen(numbers []*int) []*int {
	size := len(numbers)
	numbers = append(numbers, &size)
	return numbers
}

func main() {
	// Создаем слайс указателей с начальной емкостью 5.
	numbers := make([]*int, 0, 5)

	// ВАЖНАЯ ОШИБКА: переменная number объявлена вне цикла.
	var number int
	for range 3 {
		number++ // number = 1, 2, 3
		// Добавляем указатель на ТУ ЖЕ переменную number.
		// Все три элемента слайса указывают на одну переменную!
		numbers = append(numbers, &number)
	}
	// После цикла number = 3, все указатели в numbers ссылаются на number.

	// appendLenWrong не изменит numbers в main, так как работает с копией.
	appendLenWrong(numbers)

	// Выводим значения через указатели.
	for _, number := range numbers {
		// Все указатели ссылаются на одну переменную number = 3.
		fmt.Printf("%d ", *number) // Выведет: 3 3 3
	}
	fmt.Println()

	// Демонстрация правильного способа:
	numbers2 := make([]*int, 0, 5)
	var num2 int
	for range 3 {
		num2++
		numbers2 = append(numbers2, &num2)
	}
	// Правильно: сохраняем возвращенный слайс
	numbers2 = appendLen(numbers2)
	fmt.Print("Правильный способ: ")
	for _, n := range numbers2 {
		fmt.Printf("%d ", *n) // Выведет: 3 3 3 3
	}
	fmt.Println()
}

// Ответ: 3 3 3
//
// Объяснение:
//
// 1. Проблема с циклом:
//    var number int
//    for range 3 {
//        number++
//        numbers = append(numbers, &number)
//    }
//    - number объявлена вне цикла, это ОДНА переменная
//    - На каждой итерации number увеличивается: 1, 2, 3
//    - Но мы добавляем указатель на ТУ ЖЕ переменную
//    - После цикла все три указателя ссылаются на number = 3
//    - Поэтому вывод: 3 3 3
//
// 2. Проблема с appendLenWrong:
//    func appendLenWrong(numbers []*int) {
//        numbers = append(numbers, &size)
//    }
//    - Слайс передается по значению (копируется header)
//    - append изменяет локальную переменную numbers
//    - В main слайс numbers не изменяется
//    - Даже если бы append изменил базовый массив,
//      в main len и cap остались бы прежними
//
// 3. Как правильно исправить appendLenWrong:
//
//    ИДИОМАТИЧНЫЙ СПОСОБ: Возвращать измененный слайс
//    func appendLen(numbers []*int) []*int {
//        size := len(numbers)
//        return append(numbers, &size)
//    }
//    Вызов: numbers = appendLen(numbers)
//
//    Это стандартный паттерн в Go для функций, изменяющих слайсы.
//    Примеры из стандартной библиотеки:
//    - slice = append(slice, element)
//    - result = slices.Compact(data)
//    - filtered = slices.Delete(items, i, j)
//
//    Альтернативный способ (реже используется):
//    Можно передавать указатель на слайс *[]*int, но это менее идиоматично
//    и усложняет код без реальной выгоды.
//
// 4. Как правильно исправить цикл:
//
//    Вариант 1: Объявлять переменную внутри цикла
//    for i := range 3 {
//        number := i + 1
//        numbers = append(numbers, &number)
//    }
//
//    Вариант 2: Копировать значение перед взятием адреса
//    for range 3 {
//        number++
//        n := number
//        numbers = append(numbers, &n)
//    }
//
// Ключевые моменты:
//
// - Слайсы передаются по значению (копируется header {ptr, len, cap})
// - append может изменить базовый массив (если хватает cap)
// - append может создать новый массив (если cap недостаточно)
// - Присваивание numbers = append(...) меняет только локальную переменную
// - Чтобы вернуть изменения: возвращайте слайс или передавайте указатель
// - Переменная цикла в Go не пересоздается на каждой итерации
// - Адрес переменной цикла одинаков на всех итерациях
