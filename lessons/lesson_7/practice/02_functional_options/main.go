/*
Задание 2. Functional Options.

Дана структура Server с полями:
  Host    string
  Port    int
  Timeout int

Host обязателен (приходит в конструктор). Port и Timeout — опциональны.

Реализуй:
  - тип Option func(*Server);
  - функции WithPort(port int) Option и WithTimeout(timeout int) Option;
  - конструктор NewServer(host string, opts ...Option) *Server.

По умолчанию Port = 8080, Timeout = 30.

Подсказка: задай дефолты ДО применения опций, затем пройди циклом по opts
и примени каждую к &server.
*/

package main

import "fmt"

type Server struct {
	Host    string
	Port    int
	Timeout int
}

// TODO: тип Option

// TODO: WithPort, WithTimeout

// TODO: NewServer
func NewServer(host string /* opts ...Option */) *Server {
	return nil
}

func main() {
	s := NewServer("localhost" /*, WithPort(9090) */)
	fmt.Printf("%+v\n", s)
}
