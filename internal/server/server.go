package server

import (
	"log"

	"github.com/gofiber/fiber/v2"

	"teenyurl/internal/database"
)

type FiberServer struct {
	*fiber.App

	db database.Service
}

func New() *FiberServer {
	server := &FiberServer{
		App: fiber.New(fiber.Config{
			ServerHeader: "teenyurl",
			AppName:      "teenyurl",
		}),

		db: database.New(),
	}
	err := server.db.Init()
	if err != nil {
		log.Fatal(err)
	}
	return server
}
