package server

import (
	"github.com/gofiber/fiber/v2"

	"tft-dps-simulator/internal/database"
)

type FiberServer struct {
	*fiber.App

	db database.Service
}

func New() *FiberServer {
	server := &FiberServer{
		App: fiber.New(fiber.Config{
			ServerHeader: "tft-dps-simulator",
			AppName:      "tft-dps-simulator",
		}),

		db: database.New(),
	}

	return server
}
