package server

import (
	"fmt"
	"log"

	"tft-dps-simulator/internal/service" 

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func (s *FiberServer) RegisterFiberRoutes() {
	// Apply CORS middleware
	s.App.Use(cors.New(cors.Config{
		AllowOrigins:     "*", // Allow requests from your frontend origin
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS,PATCH",
		AllowHeaders:     "Accept,Authorization,Content-Type,X-CSRF-Token", // Ensure Content-Type is allowed
		AllowCredentials: false,                                            // Set to true if needed, adjust AllowOrigins accordingly
		MaxAge:           300,
	}))

	// Define API group v1
	apiV1 := s.App.Group("/api/v1")

	// Simulation routes
	simulationGroup := apiV1.Group("/simulation")
	
	// Use real implementation
	simulationGroup.Post("/run", s.HandleRunSimulation)
	
	// Add mock endpoint for testing
	simulationGroup.Post("/mock-run", func(c *fiber.Ctx) error {
		// Parse request just to validate it (optional)
		var req service.RunSimulationRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Cannot parse request body",
			})
		}
		
		 // Call the MockRunSimulation service method instead of duplicating response logic
		resp, err := s.simService.MockRunSimulation(req.BoardChampions)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": fmt.Sprintf("Mock simulation failed: %v", err),
			})
		}
		
		return c.Status(fiber.StatusOK).JSON(resp)
	})

	// Existing routes (keep them if needed, or move under API group)
	// s.App.Get("/", s.HelloWorldHandler)
	s.App.Get("/health", s.healthHandler)
}

// HandleRunSimulation handles requests to run the combat simulation.
func (s *FiberServer) HandleRunSimulation(c *fiber.Ctx) error {
	// 1. Parse Request Body
	var req service.RunSimulationRequest
	if err := c.BodyParser(&req); err != nil {
		log.Printf("Error parsing request body: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse request body",
		})
	}

	// 2. Basic Validation (Example: check if champions are provided)
	if len(req.BoardChampions) == 0 {
		log.Println("Validation Error: No board champions provided")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "BoardChampions array cannot be empty",
		})
	}

	// 3. Call Simulation Service
	log.Printf("Calling SimulationService with %d champions", len(req.BoardChampions))
	resp, err := s.simService.RunSimulation(req.BoardChampions)
	if err != nil {
		log.Printf("Error running simulation: %v", err)
		// Determine appropriate status code based on error type if possible
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Simulation failed: %v", err),
		})
	}

	// 4. Send Response
	log.Println("Simulation successful, sending response.")
	return c.Status(fiber.StatusOK).JSON(resp)
}

func (s *FiberServer) HelloWorldHandler(c *fiber.Ctx) error {
	resp := fiber.Map{
		"message": "Hello World",
	}

	return c.JSON(resp)
}

func (s *FiberServer) healthHandler(c *fiber.Ctx) error {
	// Assuming db is removed or handled differently if not needed
	// return c.JSON(s.db.Health())
	return c.JSON(fiber.Map{"status": "ok"}) // Simple health check if DB is removed
}
