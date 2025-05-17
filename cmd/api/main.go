package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"tft-dps-simulator/internal/core/data"
	"tft-dps-simulator/internal/server"
	"tft-dps-simulator/internal/service" 
	"github.com/gofiber/fiber/v2"
	"time"

	_ "github.com/joho/godotenv/autoload"
)

func gracefulShutdown(fiberServer *server.FiberServer, done chan bool) {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	<-ctx.Done()
	log.Println("shutting down gracefully…")

	ctxTimeout, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := fiberServer.App.ShutdownWithContext(ctxTimeout); err != nil {
		log.Printf("forced shutdown error: %v", err)
	}
	done <- true
}

func main() {
	// 1. Load Game Data
	log.Println("Loading game data...")
	dataDir := "./assets"
	fileName := "en_us_pbe.json"
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

	simService := service.NewSimulationService(tftData)
	srv := server.New(simService)

	// 1) Serve your API routes
	srv.RegisterFiberRoutes()
	// 2) Serve the React build on disk at /frontend/dist
	//    The Dockerfile copies it to /app/frontend/dist, and WORKDIR is /app
	srv.App.Static("/", "./frontend/dist")

	// 3) If you want client‐side routing support, catch-all to index.html:
	srv.App.Use(func(c *fiber.Ctx) error {
		return c.SendFile("./frontend/dist/index.html")
	})

	// 4) Start listening on $PORT
	portStr := os.Getenv("PORT")
	if portStr == "" {
		portStr = "8080"
	}
	log.Printf("Listening on :%s", portStr)

	done := make(chan bool, 1)
	go func() {
		if err := srv.Listen(":" + portStr); err != nil {
			log.Fatalf("server error: %v", err)
		}
	}()

	go gracefulShutdown(srv, done)

	<-done
	log.Println("shutdown complete")
}
