package main

import (
	"log"

	"github.com/irmadev7/tripmate-backend/internal/server"
)

func main() {
	s := server.New()
	if err := s.Run(); err != nil {
		log.Fatal(err)
	}
}
