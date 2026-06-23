package main

import "fmt"

type Closer struct {
	handlers []func() error
}

func (c *Closer) Add(handler func() error) {
	c.handlers = append(c.handlers, handler)
}

func (c *Closer) Close() error {
	for _, handler := range c.handlers {
		if err := handler(); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	closer := &Closer{}

	closer.Add(func() error {
		fmt.Println("stop worker")
		return nil
	})
	closer.Add(func() error {
		fmt.Println("close db connection")
		return nil
	})

	if err := closer.Close(); err != nil {
		fmt.Println("close error:", err)
	}
}
