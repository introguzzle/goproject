package main

import (
	"context"
	_ "github.com/lib/pq"
	"goproject/src/server"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	statusChan := make(chan string)

	wg.Add(1)
	go func() {
		defer wg.Done()
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Recovered from panic: %v", r)
			}
		}()
		server.Serve(ctx, statusChan)
	}()

	go func() {
		for status := range statusChan {
			log.Println(status)
		}
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	<-sigs
	log.Println("Received termination signal, shutting down...")

	cancel()
	wg.Wait()
	close(statusChan)

	log.Println("Server shut down successfully")
}
