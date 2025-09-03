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

	api := s.App.Group("/api")

	celulose := api.Group("/celulose")
	celulose.Post("/", s.HandleCreateLoad)
	celulose.Get("/latest", s.HandleGetLatest)
	celulose.Get("/day", s.HandleGetSummary)
	celulose.Post("/filtered", s.HandleGetFiltered)
	celulose.Patch("/:id", s.HandleUpdateLoad)
	celulose.Delete("/:id", s.HandleDeleteLoad)
}
