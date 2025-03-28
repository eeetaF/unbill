package main

import (
	"backend/init"
	"backend/server"
)

func main() {
	init.InitializeApp()
	server.Expose()
}
