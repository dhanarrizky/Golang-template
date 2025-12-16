package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/dhanarrizky/Golang-template/internal/bootstrap"
)

func main() {
	_ = godotenv.Load()

	if err := bootstrap.RunHTTPServer(); err != nil {
		log.Fatalf("app stopped: %v", err)
	}
}
