package main

import (
	"fmt"
	"votecube-id/server"
	"votecube-id/verify"
)

func main() {
	fmt.Println("Google Auth 1")

	verify.SetConfig()

	server.Start("8080", server.Dev)

}
