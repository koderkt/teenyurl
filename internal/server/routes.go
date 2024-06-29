package server

import (
	"context"
	"encoding/json"
	"fmt"
	"teenyurl/internal/types"
	"teenyurl/internal/utils"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func (s *FiberServer) RegisterFiberRoutes() {
	s.App.Get("/", s.HelloWorldHandler)

	s.App.Get("/health", s.healthHandler)
	s.App.Post("/signup", s.SignUpHandler)
	s.App.Post("/signin", s.SignInHandler)
}

func (s *FiberServer) HelloWorldHandler(c *fiber.Ctx) error {
	resp := fiber.Map{
		"message": "Hello World",
	}

	return c.JSON(resp)
}

func (s *FiberServer) SignInHandler(c *fiber.Ctx) error {

	userSignRequest := new(types.UserSignInRequest)
	err := c.BodyParser(userSignRequest)
	if err != nil {
		c.Status(400)
		c.JSON(fiber.Error{
			Code:    fiber.StatusBadRequest,
			Message: err.Error(),
		})
		return nil
	}
	// validate the user struct
	validate := validator.New()
	err = validate.Struct(userSignRequest)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": "email or password is incorrect",
		})
	}

	// Get User by email
	user, err := s.db.GetUserByEmail(userSignRequest.Email)
	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": "email or password is incorrect",
		})
	}

	// Validate password
	isValidPassword := utils.CheckPasswordHash(userSignRequest.Password, user.EncryptedPassword)

	if !isValidPassword {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  fiber.StatusUnauthorized,
			"message": "email or password is incorrect",
		})
	}

	// Generate session_id
	sessionId := uuid.NewString()

	userSession, err := json.Marshal(types.UserSession{
		Id:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  fiber.ErrInternalServerError,
			"message": "internal server error",
		})
	}

	// Store session_id and send it to client
	err = s.redisClient.Set(context.Background(), sessionId, string(userSession), 2*time.Hour).Err()
	if err != nil {
		fmt.Println(err)

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  fiber.ErrInternalServerError,
			"message": "internal server error",
		})
	}

	c.Response().Header.Set("Authorization", fmt.Sprintf("Bearer %s", sessionId))

	return c.JSON(fiber.Map{"success": true})
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
	validate := validator.New()
	err = validate.Struct(userCreationRequest)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": "bad request",
		})
	}

	// Validate password
	validate = validator.New()
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
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Error{
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
