package main

import (
	"fmt"
	"votecube-id/server"
)

func main() {
	fmt.Println("Hello, world 2")

	server.Start("8080")

}
