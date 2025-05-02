package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"tft-dps-simulator/internal/core/data" // Import data package
	"tft-dps-simulator/internal/server"
	"tft-dps-simulator/internal/service" // Import service package
	"time"
	"path/filepath"

	_ "github.com/joho/godotenv/autoload"
)

func gracefulShutdown(fiberServer *server.FiberServer, done chan bool) {
	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Listen for the interrupt signal.
	<-ctx.Done()

	log.Println("shutting down gracefully, press Ctrl+C again to force")
	stop() // Allow Ctrl+C to force shutdown

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctxTimeout, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Use the App field directly for ShutdownWithContext if FiberServer embeds *fiber.App
	// Assuming FiberServer has an 'App' field of type *fiber.App
	if err := fiberServer.App.ShutdownWithContext(ctxTimeout); err != nil {
		log.Printf("Server forced to shutdown with error: %v", err)
	}

	log.Println("Server exiting")

	// Notify the main goroutine that the shutdown is complete
	done <- true
}

func main() {
	// 1. Load Game Data
	log.Println("Loading game data...")
	dataDir := "../../assets"
	fileName := "en_us_14.1b.json"
	filePath := filepath.Join(dataDir, fileName)
	tftData, err := data.LoadSetDataFromFile(filePath, "TFTSet14")
	if err != nil {
		log.Printf("Error loading set data: %v\n", err)
		os.Exit(1)
	}
	data.InitializeChampions(tftData)

	data.InitializeTraits(tftData)
	data.InitializeSetActiveAugments(tftData, filePath)

	data.InitializeSetActiveItems(tftData, filePath)


	// 2. Initialize Services
	simService := service.NewSimulationService(tftData)

	// 3. Initialize Server with Services
	server := server.New(simService) // Pass simService to New

	server.RegisterFiberRoutes()

	// Create a done channel to signal when the shutdown is complete
	done := make(chan bool, 1)

	go func() {
		portStr := os.Getenv("PORT")
		if portStr == "" {
			portStr = "8080" // Default port
		}
		port, err := strconv.Atoi(portStr)
		if err != nil {
			log.Fatalf("Invalid PORT environment variable: %s", portStr)
		}

		err = server.Listen(fmt.Sprintf(":%d", port))
		if err != nil {
			// Use log.Fatalf to exit if server fails to start
			log.Fatalf("HTTP server error: %s", err)
		}
	}()

	// Run graceful shutdown in a separate goroutine
	go gracefulShutdown(server, done)

	// Wait for the graceful shutdown to complete
	<-done
	log.Println("Graceful shutdown complete.")
}
