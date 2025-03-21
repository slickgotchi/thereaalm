package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"thereaalm/network"
	"thereaalm/world"
)

func main() {
	log.Println("Starting The Reaalm...")

	// Create and run the world manager
	worldManager := world.NewWorldManager(0) // Use all CPU cores
	worldManager.Run()

	// start the api server
	network.StartAPIServer(worldManager, "8080")

	// Create a channel to listen for interrupt signals (e.g., Ctrl+C)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Block until an interrupt signal is received
	<-sigChan
	log.Println("Received shutdown signal, exiting...")
}