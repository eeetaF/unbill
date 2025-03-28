package main

import (
	"backend/initialization"
	"backend/server"
)

func main() {
	initialization.InitializeApp()
	server.Expose()
}
