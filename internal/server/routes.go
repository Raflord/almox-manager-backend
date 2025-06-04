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
	celulose.Post("/", s.HandleCreate)
	celulose.Get("/latest", s.HandleGetLatest)
	celulose.Get("/day", s.HandleGetDay)
	celulose.Post("/filtered", s.HandleFiltered)
	celulose.Put("/:id", s.HandleUpdate)
	celulose.Delete("/:id", s.HandleDelete)

	// insumos := api.Group("/insumos")
	// insumos.Post("/")
	// insumos.Get("/latest")
	// insumos.Get("/day")
	// insumos.Post("/filtered")
	// insumos.Put("/:id")
	// insumos.Delete("/:id")
}
