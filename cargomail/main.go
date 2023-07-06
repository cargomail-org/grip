package main

import (
	"cargomail/provider"
	"log"
)

func main() {
	err := provider.Start()
	if err != nil {
		log.Fatalf("provider service error: %v", err)
	}
	log.Print("provider service shutdown gracefully")
}
