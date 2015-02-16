package main

import (
	"./server"
	"fmt"
	"net/http"
)

func main() {
	server := server.New()

	server.Method("hello", func(p []interface{}) (interface{}, error) {
		return "result", nil
	})

	http.HandleFunc("/websocket", server.Handler)
	http.ListenAndServe(":9000", nil)
}
