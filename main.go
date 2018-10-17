package main

import (
	"fmt"
	"votecube-id/server"
)

func main() {
	fmt.Println("Hello, world 3")

	server.Start("8080", server.Dev)

}
