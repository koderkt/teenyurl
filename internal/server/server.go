package server

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"

	"teenyurl/internal/database"
)

type FiberServer struct {
	*fiber.App
	redisClient *redis.Client
	db database.Service
}

func New() *FiberServer {
	server := &FiberServer{
		App: fiber.New(fiber.Config{
			ServerHeader: "teenyurl",
			AppName:      "teenyurl",
		}),
		redisClient: database.CreateRedisConnection(),
		db: database.New(),
	}
	err := server.db.Init()
	if err != nil {
		log.Fatal(err)
	}
	return server
}
