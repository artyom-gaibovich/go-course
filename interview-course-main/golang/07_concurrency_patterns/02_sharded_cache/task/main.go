// Задача: Напишите реализацию InMemory кэша

package main

type Cache interface {
	Set(k string, v string)
	Get(k string) (string, bool)
}
