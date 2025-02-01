package main

import (
	"log"

	"github.com/developerc/project_gophermart/internal/server"
)

func main() {
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
