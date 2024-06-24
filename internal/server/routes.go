package server

import (
	"teenyurl/internal/types"
	"time"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

func (s *FiberServer) RegisterFiberRoutes() {
	s.App.Get("/", s.HelloWorldHandler)

	s.App.Get("/health", s.healthHandler)
	s.App.Post("/signup", s.SignUpHandler)
}

func (s *FiberServer) HelloWorldHandler(c *fiber.Ctx) error {
	resp := fiber.Map{
		"message": "Hello World",
	}

	return c.JSON(resp)
}

func (s *FiberServer) SignUpHandler(c *fiber.Ctx) error {
	userCreationRequest := new(types.CreateUserRequest)
	err := c.BodyParser(userCreationRequest)
	if err != nil {
		c.Status(400)
		c.JSON(fiber.Error{
			Code:    fiber.StatusBadRequest,
			Message: err.Error(),
		})
		return nil
	}
	encpw, err := bcrypt.GenerateFromPassword([]byte(userCreationRequest.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user := types.User{
		FirstName:         userCreationRequest.FirstName,
		LastName:          userCreationRequest.LastName,
		Email:             userCreationRequest.Email,
		EncryptedPassword: string(encpw),
		CreatedAt:         time.Now(),
	}
	err = s.db.CreateUser(&user)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Error{
		Code:    fiber.StatusAccepted,
		Message: "user created",
	})
}
func (s *FiberServer) healthHandler(c *fiber.Ctx) error {
	return c.JSON(s.db.Health())
}
