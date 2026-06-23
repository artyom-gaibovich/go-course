package main

import "fmt"

type Logger struct {
	name  string
	level string
}

func NewLogger() *Logger {
	return &Logger{}
}

func (l *Logger) WithName(name string) *Logger {
	l.name = name
	return l
}

func (l *Logger) WithLevel(level string) *Logger {
	l.level = level
	return l
}

func main() {
	logger := NewLogger().
		WithName("storage").
		WithLevel("INFO")

	fmt.Printf("%+v\n", *logger)
}
