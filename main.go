package main

import (
	"log"

	"github.com/irmadev7/tripmate-backend/internal/server"
)

func main() {
	s, err := server.New()
	if err != nil {
		log.Fatal(err)
	}
	if err := s.Run(); err != nil {
		log.Fatal(err)
	}
}
