package main

import (
	_ "proiect_web/mongodb"
	"proiect_web/server"
)

func main() {
	server.RunServer("8080")
}
