/*
Задание 4. Паттерн Closer.

Реализуй компонент Closer, который собирает обработчики освобождения ресурсов
и выполняет их в одном месте.

Требования:
  - метод Add(name string, handler func() error) — регистрирует обработчик;
  - метод Close() error — вызывает обработчики В ОБРАТНОМ порядке регистрации
    (как defer / как LIFO разворачивание ресурсов);
    при первой ошибке прекращает выполнение и возвращает её.

Программа ниже должна напечатать close в порядке: server, db, worker.

Подсказка: ресурсы обычно освобождают в порядке, обратном захвату — поэтому Close
идёт с конца слайса к началу.
*/

package main

import "fmt"

type Closer struct {
	// TODO: поле для хранения обработчиков
}

// TODO: Add(name string, handler func() error)

// TODO: Close() error

func main() {
	c := &Closer{}
	c.Add("worker", func() error { fmt.Println("close worker"); return nil })
	c.Add("db", func() error { fmt.Println("close db"); return nil })
	c.Add("server", func() error { fmt.Println("close server"); return nil })

	if err := c.Close(); err != nil {
		fmt.Println("error:", err)
	}
}
