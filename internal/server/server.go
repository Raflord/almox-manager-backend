package server

import (
	"github.com/gofiber/fiber/v2"

	"almox-manager-backend/internal/database"
)

type FiberServer struct {
	*fiber.App

	db database.Service
}

func New() *FiberServer {
	server := &FiberServer{
		App: fiber.New(fiber.Config{
			ServerHeader: "almox-manager-backend",
			AppName:      "almox-manager-backend",
		}),

		db: database.New(),
	}

	return server
}
