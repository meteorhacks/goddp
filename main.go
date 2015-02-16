package main

import (
	"./server"
)

func main() {
	server := server.New()
	server.Method("hello", func(p []interface{}) (interface{}, error) {
		return "result", nil
	})

	server.Listen(":9000")
}
