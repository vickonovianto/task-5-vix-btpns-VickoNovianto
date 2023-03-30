package main

import (
	"log"
	"sync"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Panic(err)
	}
	server := InitServer()
	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		server.Run()
	}()

	wg.Wait()
}
