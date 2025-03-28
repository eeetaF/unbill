package main

import (
	runtime_data "backend/runtime-data"
	"backend/server"
)

func main() {
	runtime_data.InitializeApp()
	server.Expose()
}
