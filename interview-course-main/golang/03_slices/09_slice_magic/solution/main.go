package main

import "fmt"

func main() {
	// Создаем слайс с len=10, cap=10
	original := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	fmt.Printf("original: len=%d cap=%d %v\n", len(original), cap(original), original)

	// Создаем срез от слайса [2:5] - элементы с индексами 2, 3, 4
	slice1 := original[2:5]
	fmt.Printf("slice1 [2:5]: len=%d cap=%d %v\n", len(slice1), cap(slice1), slice1)

	// Создаем еще один срез [3:7] - элементы с индексами 3, 4, 5, 6
	slice2 := original[3:7]
	fmt.Printf("slice2 [3:7]: len=%d cap=%d %v\n", len(slice2), cap(slice2), slice2)

	// Модифицируем первый элемент slice1
	slice1[0] = 999
	fmt.Printf("\nПосле slice1[0]=999:\n")
	fmt.Printf("original: %v\n", original)
	fmt.Printf("slice1: %v\n", slice1)
	fmt.Printf("slice2: %v\n", slice2)

	// Добавляем элементы к slice1
	slice1 = append(slice1, 100, 200)
	fmt.Printf("\nПосле append(slice1, 100, 200):\n")
	fmt.Printf("slice1: len=%d cap=%d %v\n", len(slice1), cap(slice1), slice1)
	fmt.Printf("original: %v\n", original)
	fmt.Printf("slice2: %v\n", slice2)
}

// Вывод:
// original: len=10 cap=10 [0 1 2 3 4 5 6 7 8 9]
// slice1 [2:5]: len=3 cap=8 [2 3 4]
// slice2 [3:7]: len=4 cap=7 [3 4 5 6]
//
// После slice1[0]=999:
// original: [0 1 999 3 4 5 6 7 8 9]
// slice1: [999 3 4]
// slice2: [999 4 5 6]
//
// После append(slice1, 100, 200):
// slice1: len=5 cap=8 [999 3 4 100 200]
// original: [0 1 999 3 4 100 200 7 8 9]
// slice2: [999 4 100 200]

// ============================================================================
// ОБЪЯСНЕНИЕ: МАГИЯ CAPACITY ПРИ СОЗДАНИИ СРЕЗОВ
// ============================================================================
//
// 1. ВНУТРЕННЯЯ СТРУКТУРА СЛАЙСА
//
// Слайс в Go это структура из 3 полей:
//   type slice struct {
//       array unsafe.Pointer  // Указатель на базовый массив
//       len   int              // Длина (количество доступных элементов)
//       cap   int              // Емкость (до конца базового массива)
//   }
//
// 2. КАК РАБОТАЕТ SLICE EXPRESSION [low:high]
//
// Синтаксис: original[low:high]
// - low  - начальный индекс (включительно)
// - high - конечный индекс (НЕ включительно)
//
// Что происходит:
//   newSlice := original[low:high]
//   newSlice.array = &original.array[low]  // Указатель сдвигается на low
//   newSlice.len   = high - low            // Новая длина
//   newSlice.cap   = original.cap - low    // МАГИЯ ЗДЕСЬ!
//
// КЛЮЧЕВОЙ МОМЕНТ:
// Capacity нового слайса = capacity оригинала МИНУС смещение low!
// Это потому что capacity это расстояние от ТЕКУЩЕЙ позиции до конца массива.
//
// 3. ПРИМЕР С original[2:5]
//
// original = [0 1 2 3 4 5 6 7 8 9]
//             ^               ^
//            ptr            cap=10
//
// slice1 := original[2:5]
//             ↓
// original = [0 1 2 3 4 5 6 7 8 9]
//                 ^ ^ ^       ^
//                 │ │ │       │
//                 │ │ │       └─ cap boundary (original[9])
//                 │ │ └─ high (not included)
//                 │ └─ elements in slice1
//                 └─ low (new ptr position)
//
// Расчет:
//   slice1.array = &original[2]     // Указатель на элемент 2
//   slice1.len   = 5 - 2 = 3        // Элементы: 2, 3, 4
//   slice1.cap   = 10 - 2 = 8       // От позиции 2 до конца массива
//
// Визуализация capacity:
//   original[2:5] может "видеть" элементы 2, 3, 4, 5, 6, 7, 8, 9
//   Всего 8 элементов - это и есть capacity!
//
// 4. ПОЧЕМУ CAPACITY = cap - low?
//
// Потому что slice сохраняет указатель на НАЧАЛО среза, а capacity это
// расстояние от этого указателя до конца базового массива.
//
// Диаграмма памяти:
//   Базовый массив: [0][1][2][3][4][5][6][7][8][9]
//                          ↑                    ↑
//                       slice1.ptr         end of array
//                          |<--- cap = 8 ----->|
//                          |<- len=3 ->|
//
// 5. ПОЧЕМУ ВСЕ СЛАЙСЫ РАЗДЕЛЯЮТ БАЗОВЫЙ МАССИВ
//
// original, slice1, slice2 - это ТРИ РАЗНЫХ слайса (3 структуры),
// но они ВСЕ указывают на ОДИН И ТОТ ЖЕ базовый массив!
//
//   Memory:
//   ┌────────────────────────────────────────┐
//   │ Array: [0][1][999][3][4][5][6][7][8][9]│
//   └────────────────────────────────────────┘
//       ↑           ↑     ↑
//       │           │     │
//   original    slice1  slice2
//   ptr=&[0]    ptr=&[2] ptr=&[3]
//   len=10      len=3    len=4
//   cap=10      cap=8    cap=7
//
// Поэтому:
// - slice1[0] = 999 меняет original[2]
// - Это же изменение видно в slice2[0] (так как slice2 начинается с original[3])
//
// 6. ЧТО ПРОИСХОДИТ ПРИ APPEND
//
// slice1 = append(slice1, 100, 200)
//
// До append:
//   slice1: len=3, cap=8, [999, 3, 4]
//   Можем добавить еще 5 элементов без реаллокации (cap - len = 8 - 3 = 5)
//
// После добавления 2 элементов:
//   slice1: len=5, cap=8, [999, 3, 4, 100, 200]
//   Реаллокации НЕ было, так как 5 <= 8
//
// Куда записались 100 и 200?
//   В позиции original[5] и original[6]!
//   original = [0, 1, 999, 3, 4, 100, 200, 7, 8, 9]
//                                ^^^  ^^^
//
// Почему это влияет на slice2?
//   slice2 = original[3:7] = [original[3], original[4], original[5], original[6]]
//   После append: [999, 4, 100, 200]
//                       ^^^  ^^^
//
// 7. FULL SLICE EXPRESSION [low:high:max]
//
// Можно явно указать capacity через third index:
//   slice3 := original[2:5:5]
//   // len = 5-2 = 3
//   // cap = 5-2 = 3 (вместо 8!)
//
// Это защищает от случайного изменения элементов за пределами high.
//
// 8. КОГДА ПРОИСХОДИТ РЕАЛЛОКАЦИЯ
//
// Если бы мы добавили БОЛЬШЕ элементов:
//   slice1 = append(slice1, 100, 200, 300, 400, 500, 600)
//   // Нужно добавить 6 элементов, но cap-len = 5
//
// Тогда:
//   1. Go создаст новый массив большего размера
//   2. Скопирует данные из старого массива
//   3. slice1 будет указывать на НОВЫЙ массив
//   4. original и slice2 остаются на СТАРОМ массиве
//   5. Дальнейшие изменения slice1 НЕ влияют на original/slice2
//
// 9. ПРАКТИЧЕСКИЕ ВЫВОДЫ
//
// - Срез от слайса РАЗДЕЛЯЕТ базовый массив с оригиналом
// - Capacity среза = capacity оригинала минус смещение
// - Изменения в одном срезе видны в других (если они перекрываются)
// - append может изменить "невидимые" элементы базового массива
// - Используйте full slice expression [low:high:max] для контроля capacity
// - При необходимости копируйте: newSlice := make([]T, len(old)); copy(newSlice, old)
//
// 10. ТИПИЧНЫЕ ОШИБКИ
//
// Ошибка 1: Ожидание независимости срезов
//   s1 := original[0:5]
//   s2 := original[5:10]
//   s1[0] = 999           // OK, не влияет на s2
//   s1 = append(s1, 100)  // Запишет в original[5] - ВЛИЯЕТ НА s2!
//
// Ошибка 2: Возврат среза от локального слайса
//   func bad() []int {
//       data := []int{1, 2, 3, 4, 5}
//       return data[2:4]  // Возвращаем срез, но базовый массив может быть GC
//   }
//   // Лучше копировать:
//   result := make([]int, 2)
//   copy(result, data[2:4])
//   return result
//
// Ошибка 3: Итерация по изменяемому слайсу
//   for i := range slice {
//       slice = append(slice, i)  // Бесконечный цикл!
//   }
//   // Правильно:
//   n := len(slice)
//   for i := 0; i < n; i++ { ... }
