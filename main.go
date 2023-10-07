package main

import "video/internal/app/converter"

func main() {
	c := &converter.App{}
	err := c.Register()
	if err != nil {
		panic(err)
	}
	err = c.Run()
	if err != nil {
		c.Resolve()
	}
}
