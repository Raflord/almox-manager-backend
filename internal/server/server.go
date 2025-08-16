package server

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/joho/godotenv/autoload"

	"almox-manager-backend/internal/database"
)

type FiberServer struct {
	*fiber.App
	db *database.Queries
}

type service struct {
	db *pgxpool.Pool
}

var (
	dbname     = os.Getenv("DB_DATABASE")
	password   = os.Getenv("DB_PASSWORD")
	username   = os.Getenv("DB_USERNAME")
	port       = os.Getenv("DB_PORT")
	host       = os.Getenv("DB_HOST")
	pgOnce     sync.Once
	dbInstance *service
)

func New() *FiberServer {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", username, password, host, port, dbname)
	pgOnce.Do(func() {
		db, err := pgxpool.New(context.Background(), connStr)
		if err != nil {
			log.Fatal(err)
		}

		dbInstance = &service{db}
	})

	server := &FiberServer{
		App: fiber.New(fiber.Config{
			ServerHeader: "almox-manager-backend",
			AppName:      "almox-manager-backend",
		}),

		db: database.New(dbInstance.db),
	}

	return server
}
