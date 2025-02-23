package server

import (
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func (s *FiberServer) RegisterFiberRoutes() {
	// Apply CORS middleware
	s.App.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS,PATCH",
		AllowHeaders:     "Accept,Authorization,Content-Type",
		AllowCredentials: false, // credentials require explicit origins
		MaxAge:           300,
	}))

	s.App.Get("/api/cellulose/latest", s.HandleGetLatestCell)
	s.App.Get("/api/cellulose/day", s.HandleGetDayCell)
	s.App.Post("/api/cellulose/filtered", s.HandleGetFiltered)
	s.App.Post("/api/cellulose", s.HandlePostCellulose)
	s.App.Put("/api/cellulose", s.HandlePutCellulose)
	s.App.Delete("/api/cellulose", s.HandleDeleteCellulose)
	s.App.Get("/api/cellulose", s.handleGetCellulose)
}

func (s *FiberServer) handleGetCellulose(c *fiber.Ctx) error {
	panic("crash")
}
