package server

import (
	"net/mail"
	"teenyurl/internal/types"
	"teenyurl/internal/utils"
	"time"

	"github.com/go-playground/validator/v10"
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

	if len(userCreationRequest.Email) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": "email not entered",
		})
	}

	if len(userCreationRequest.FirstName) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": "first_name not entered",
		})
	}

	if len(userCreationRequest.LastName) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": "last_name not entered",
		})
	}
	// check if email is in approriate format
	_, err = mail.ParseAddress(userCreationRequest.Email)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": "invalid email format",
		})
	}

	// Validate password
	validate := validator.New()
	validate.RegisterValidation("password", utils.PasswordValidator)
	if err := validate.Var(userCreationRequest.Password, "required,password"); err != nil {
		return c.JSON(fiber.Error{
			Code:    fiber.StatusBadRequest,
			Message: "invalid password",
		})
	}
	// Encrypt password before storing
	encpw, err := bcrypt.GenerateFromPassword([]byte(userCreationRequest.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.JSON(fiber.Error{
			Code:    fiber.StatusInternalServerError,
			Message: "server error",
		})
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
		return c.JSON(fiber.Error{
			Code:    fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	return c.JSON(fiber.Error{
		Code:    fiber.StatusAccepted,
		Message: "user created",
	})
}
func (s *FiberServer) healthHandler(c *fiber.Ctx) error {
	return c.JSON(s.db.Health())
}
