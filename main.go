package main

//go:generate sqlboiler --wipe crdb

import (
	"fmt"
	"votecube-id/db"
	"votecube-id/server"
	"votecube-id/verify"

	_ "github.com/lib/pq"
)

func main() {
	fmt.Println("Google Auth 1")

	verify.SetConfig()

	dBase := db.SetupDb()

	server.Start("8080", server.Dev, dBase)

}
