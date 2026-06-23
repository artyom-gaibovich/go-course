/*
Домашнее задание. Упаковка игрового объекта в 64 байта.

Представь, что ты пишешь игровой движок. Объектов-игроков очень много, их нужно
обсчитывать максимально кэш-френдли, поэтому данные одного игрока должны помещаться
ровно в одну кэш-линию — 64 байта.

Дано (логические атрибуты игрока):
  - Name      — имя, до 42 символов ASCII (нельзя хранить как обычный string:
                string = указатель + длина, это прыжок по памяти);
  - X, Y      — координаты, int32 каждая;
  - Gold      — золото, до 16_777_215 (влезает в 3 байта / 24 бита);
  - Health    — здоровье 0..100 (хватает 7 бит);
  - Mana      — мана 0..100 (хватает 7 бит);
  - Level     — уровень 0..63 (6 бит);
  - Class     — класс 0..15 (4 бита: воин/маг/...);
  - Experience — опыт, до 4 млрд (uint32).

Требования:
  1. unsafe.Sizeof(GamePlayer{}) ДОЛЖЕН быть равен 64.
  2. Имя хранить как массив байт [N]byte (своё «отображение» строки), а не string.
  3. Поля, которые сами по себе не байтовой ширины (Health, Mana, Level, Class),
     упаковать битовыми операциями в одно числовое поле.
  4. Реализовать функциональные опции (WithGold, WithHealth, WithMana, WithLevel,
     WithClass, WithExperience, WithPosition) и конструктор NewGamePlayer(name, opts...).
  5. Реализовать геттеры (GetName, GetHealth, GetMana, GetLevel, GetClass, ...),
     которые распаковывают значения обратно.

Подсказки:
  - вспомни тему ВЫРАВНИВАНИЯ: располагай широкие поля раньше узких, чтобы не
    появлялся padding и сумма уложилась в 64 байта;
  - вспомни БИТОВЫЕ ОПЕРАЦИИ: упаковка value << shift | acc, распаковка
    (acc >> shift) & mask;
  - проверь итог через unsafe.Sizeof.
*/

package main

import (
	"fmt"
	"unsafe"
)

type GamePlayer struct {
	// TODO: спроектируй раскладку ровно на 64 байта
}

type Option func(*GamePlayer)

// TODO: WithGold, WithHealth, WithMana, WithLevel, WithClass, WithExperience, WithPosition

func NewGamePlayer(name string, opts ...Option) *GamePlayer {
	// TODO: создать игрока, записать имя в байтовый массив, применить опции
	return nil
}

// TODO: геттеры

func main() {
	fmt.Println("size:", unsafe.Sizeof(GamePlayer{}))
}
