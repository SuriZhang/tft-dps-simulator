package server

import (
	"log"
	"tft-dps-simulator/internal/service" 
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

type FiberServer struct {
	*fiber.App
	simService *service.SimulationService // Add SimulationService field
}

func New(simService *service.SimulationService) *FiberServer { // Accept simService
	app := fiber.New(fiber.Config{
		ServerHeader: "tft-dps-simulator",
		AppName:      "tft-dps-simulator",
	})

	// Add logger middleware
	app.Use(logger.New())

	server := &FiberServer{
		App:        app,
		simService: simService, // Store the service instance
	}

	return server
}

func (s *FiberServer) ShutdownWithContext(ctx fiber.Ctx) error {
	log.Println("Attempting to shutdown Fiber server...")
	return s.App.Shutdown()
}

func (s *FiberServer) Listen(addr string) error {
	log.Printf("Starting server on %s", addr)
	return s.App.Listen(addr)
}
