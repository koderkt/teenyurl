package server

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"strconv"

	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/koderkt/teenyurl/internal/types"
	"github.com/koderkt/teenyurl/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

func (s *FiberServer) RegisterFiberRoutes() {
	s.App.Get("/", s.HelloWorldHandler)
	s.App.Get("/health", s.healthHandler)
	s.App.Post("/signup", s.SignUpHandler)
	s.App.Post("/signin", s.SignInHandler)
	s.App.Post("/signout", s.SignOutHandler)
	s.App.Post("/links", s.CreateShortURLHandler)
	s.App.Get("/links", s.GetLinksHandler)
	s.App.Get("/:shortCode", s.ShortURLHandler)

	s.App.Post("/:shortCode", s.EditLongURLHandler)
	s.App.Post("/:shortCode/:val", s.EnableDisbaleURLHandler)
	s.App.Get("/analytics/:shortCode", s.AnalyticsHandler)
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
			"success": false,

			"message": "email or password is incorrect",
		})
	}

	// Get User by email
	user, err := s.db.GetUserByEmail(userSignRequest.Email)
	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,

			"message": "email or password is incorrect",
		})
	}

	// Validate password
	isValidPassword := utils.CheckPasswordHash(userSignRequest.Password, user.EncryptedPassword)

	if !isValidPassword {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,

			"message": "email or password is incorrect",
		})
	}

	// Generate session_id
	sessionId := uuid.NewString()
	userSession := types.UserSession{
		Id:       user.ID,
		UserName: user.UserName,
		Email:    user.Email,
	}
	userSessionBytes, err := json.Marshal(userSession)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,

			"message": "internal server error",
		})
	}

	// Store session_id and send it to client
	err = s.redisClient.Set(context.Background(), "session:"+sessionId, string(userSessionBytes), 2*time.Hour).Err()
	if err != nil {
		fmt.Println(err)

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "internal server error",
		})
	}

	c.Response().Header.Set("Authorization", fmt.Sprintf("Bearer %s", sessionId))

	return c.JSON(fiber.Map{"success": true,
		"user": userSession,
	})
}

func (s *FiberServer) SignOutHandler(c *fiber.Ctx) error {
	sessionHeader := c.Get("Authorization")
	if sessionHeader == "" || len(sessionHeader) < 8 || sessionHeader[:7] != "Bearer " {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "invalid session header"})
	}
	// get the session id
	sessionId := sessionHeader[7:]
	_, err := s.GetSession("session:" + sessionId)
	if err != nil {
		c.SendStatus(401)
		return c.JSON(fiber.Map{"message": "unauthorized"})
	}
	err = s.redisClient.Del(context.Background(), "session:"+sessionId).Err()
	if err != nil {
		c.SendStatus(401)
		return c.JSON(fiber.Map{"message": "internal server error"})
	}
	return c.Status(200).JSON(fiber.Map{
		"message": "logout successful",
	})
}
func (s *FiberServer) SignUpHandler(c *fiber.Ctx) error {
	userCreationRequest := new(types.CreateUserRequest)

	err := c.BodyParser(userCreationRequest)
	if err != nil {
		log.Printf("%v | Parsing CreateUserRequest | %s", time.Now().Local(), err.Error())

		c.Status(400)
		return c.JSON(fiber.Error{
			Code:    fiber.StatusBadRequest,
			Message: err.Error(),
		})

	}
	validate := validator.New()
	err = validate.Struct(userCreationRequest)

	if err != nil {
		log.Printf("%v | | %s", time.Now().Local(), err.Error())

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": "bad request",
		})
	}

	// Validate password
	validate = validator.New()
	validate.RegisterValidation("password", utils.PasswordValidator)
	if err := validate.Var(userCreationRequest.Password, "required,password"); err != nil {
		log.Printf("%v | Validating password |%s", time.Now().Local(), err.Error())

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Error{
			Code:    fiber.StatusBadRequest,
			Message: "invalid password",
		})
	}
	// Encrypt password before storing
	encpw, err := bcrypt.GenerateFromPassword([]byte(userCreationRequest.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("%v | %s", time.Now().Local(), err.Error())

		return c.JSON(fiber.Error{
			Code:    fiber.StatusInternalServerError,
			Message: "server error",
		})
	}

	user := types.User{
		UserName:          userCreationRequest.UserName,
		Email:             userCreationRequest.Email,
		EncryptedPassword: string(encpw),
		CreatedAt:         time.Now(),
	}
	err = s.db.CreateUser(&user)
	if err != nil {
		log.Printf("%v | %s", time.Now().Local(), err.Error())

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

func (s *FiberServer) CreateShortURLHandler(c *fiber.Ctx) error {
	sessionHeader := c.Get("Authorization")
	// ensure the session header is not empty and in the correct format
	if sessionHeader == "" || len(sessionHeader) < 8 || sessionHeader[:7] != "Bearer " {
		return c.JSON(fiber.Map{"error": "invalid session header"})
	}
	// get the session id
	sessionId := sessionHeader[7:]
	userSession, err := s.GetSession("session:" + sessionId)
	if err != nil {
		c.SendStatus(401)
		return c.JSON(fiber.Map{"message": "You are not logged in..."})
	}
	longURLRequst := new(types.ShortenRequest)

	err = c.BodyParser(longURLRequst)
	if err != nil {
		c.SendStatus(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{"message": err.Error()})
	}
	parsedURL, err := url.ParseRequestURI(longURLRequst.LongUrl)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid url"})
	}

	// Check if the scheme (protocol) and host (domain) are valid
	if parsedURL.Scheme == "" || parsedURL.Host == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "invalid url"})
	}

	// Validate whether the link is valid
	genaratedShortCode := utils.GenerateShortCode(6)

	_, err = s.db.GetLink(genaratedShortCode)

	if err != sql.ErrNoRows {
		c.SendStatus(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{"message": err.Error()})
	}

	link := &types.Link{
		OriginalURL: longURLRequst.LongUrl,
		ShortURL:    genaratedShortCode,
		UserId:      userSession.Id,
	}

	err = s.db.CreateShortURL(link)
	if err != nil {
		log.Printf("%v | %s", time.Now(), err.Error())
		return c.JSON(fiber.Map{"error": "failed to create short url"})
	}

	recordFromDB, err := s.db.GetLink(link.ShortURL)

	if err != nil {
		c.SendStatus(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{"message": "failed to create short url"})
	}
	responseData := types.CreateShortURLResponse{
		ShortURL:    string(c.Request().Host()) + "/" + recordFromDB.ShortURL,
		OriginalURL: recordFromDB.OriginalURL,
		LinkId:      recordFromDB.Id,
	}
	return c.Status(fiber.StatusAccepted).JSON(responseData)
}

func (s *FiberServer) healthHandler(c *fiber.Ctx) error {
	return c.JSON(s.db.Health())
}

func (s *FiberServer) GetSession(session string) (*types.UserSession, error) {
	data, err := s.redisClient.Get(context.Background(), session).Result()
	if err != nil {
		log.Printf("%v | %s", time.Now().Local(), err.Error())

		return nil, err
	}

	var userSession types.UserSession
	err = json.Unmarshal([]byte(data), &userSession)
	if err != nil {
		log.Printf("%v | %s", time.Now().Local(), err.Error())
		return nil, err
	}

	return &userSession, nil

}

func (s *FiberServer) ShortURLHandler(c *fiber.Ctx) error {
	fmt.Println("hello")
	// sessionHeader := c.Get("Authorizati  on")

	// // ensure the session header is not empty and in the correct format
	// if sessionHeader == "" || len(sessionHeader) < 8 || sessionHeader[:7] != "Bearer " {
	// 	return c.JSON(fiber.Map{"error": "invalid session header"})
	// }
	// // get the session id
	// sessionId := sessionHeader[7:]
	// _, err := s.GetSession(sessionId)
	// if err != nil {
	// 	c.SendStatus(401)
	// 	return c.JSON(fiber.Map{"message": "You are not logged in..."})
	// }
	shortCode := c.Params("shortCode")
	link, err := s.db.GetLink(shortCode)
	fmt.Println("Before....")

	if err != nil {

		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "link not found",
		})
	}
	if !link.IsEnabled {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "link is disabled at the moment",
		})
	}
	analyticsData := types.Clicks{
		ShortCode:  shortCode,
		DeviceType: "Unknown",
		Location:   "Unknown",
	}
	fmt.Println("InsertAnalytics....")
	err = s.db.InsertAnalytics(&analyticsData)
	if err != nil {
		log.Printf("%v | %s", time.Now(), err.Error())

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to record analytics",
		})
	}
	return c.Redirect(link.OriginalURL, fiber.StatusPermanentRedirect)
}

func (s *FiberServer) AnalyticsHandler(c *fiber.Ctx) error {
	sessionHeader := c.Get("Authorization")

	// ensure the session header is not empty and in the correct format
	if sessionHeader == "" || len(sessionHeader) < 8 || sessionHeader[:7] != "Bearer " {
		return c.JSON(fiber.Map{"error": "invalid session header"})
	}
	// get the session id
	sessionId := sessionHeader[7:]
	_, err := s.GetSession(sessionId)
	if err != nil {
		c.SendStatus(401)
		return c.JSON(fiber.Map{"message": "You are not logged in..."})
	}
	shortCode := c.Params("shortCode")

	clicks, err := s.db.GetAnalystics(shortCode)
	if err != nil {
		if err == sql.ErrNoRows {
			c.SendStatus(fiber.StatusAccepted)
			return err
		} else {
			log.Printf("%v | %s", time.Now(), err.Error())

			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "something went wrong"})
		}
	}

	return c.Status(fiber.StatusAccepted).JSON(clicks)

}

func (s *FiberServer) GetLinksHandler(c *fiber.Ctx) error {
	sessionHeader := c.Get("Authorization")

	// ensure the session header is not empty and in the correct format
	if sessionHeader == "" || len(sessionHeader) < 8 || sessionHeader[:7] != "Bearer " {
		return c.Status(400).JSON(fiber.Map{"error": "invalid session header"})
	}
	// get the session id
	sessionId := sessionHeader[7:]

	user, err := s.GetSession("session:" + sessionId)
	if err != nil {
		log.Printf("%v | %s", time.Now(), err.Error())

		c.SendStatus(401)
		return c.JSON(fiber.Map{"message": "You are not logged in..."})
	}

	links, err := s.db.GetLinks(user.Id)
	linksResponse := []types.LinkResponse{}
	if err != nil {
		if err == sql.ErrNoRows {
			c.SendStatus(fiber.StatusAccepted)
			return err
		} else {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "something went wrong"})
		}
	}
	for _, link := range *links {
		var linkResponse types.LinkResponse
		clicks, err := s.db.GetNumberOfClicks(link.ShortURL)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "something went wrong"})
		}
		linkResponse.ShortURL = string(c.Request().Host()) + "/" + link.ShortURL
		linkResponse.OriginalURL = link.OriginalURL
		linkResponse.CreatedAt = link.CreatedAt
		linkResponse.Clicks = clicks
		linkResponse.IsEnabled = link.IsEnabled
		linkResponse.Id = link.Id

		linksResponse = append(linksResponse, linkResponse)
	}

	return c.Status(fiber.StatusAccepted).JSON(linksResponse)
}

func (s *FiberServer) EditLongURLHandler(c *fiber.Ctx) error {
	sessionHeader := c.Get("Authorization")

	longURL := types.CreateShortURLResponse{}

	err := c.BodyParser(&longURL)
	if err != nil {
		log.Printf("%v | %s", time.Now(), err.Error())

		return c.Status(400).JSON(fiber.Map{"error": "bad request"})
	}
	// ensure the session header is not empty and in the correct format
	if sessionHeader == "" || len(sessionHeader) < 8 || sessionHeader[:7] != "Bearer " {
		return c.Status(400).JSON(fiber.Map{"error": "invalid session header"})
	}
	// get the session id
	sessionId := sessionHeader[7:]

	_, err = s.GetSession("session:" + sessionId)
	if err != nil {
		log.Printf("%v | %s", time.Now(), err.Error())

		c.SendStatus(401)
		return c.JSON(fiber.Map{"message": "You are not logged in..."})
	}

	shortCode := c.Params("shortCode")

	link, err := s.db.GetLink(shortCode)
	linksResponse := []types.LinkResponse{}
	if err != nil {
		if err == sql.ErrNoRows {
			c.SendStatus(fiber.StatusAccepted)
			return err
		} else {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "something went wrong"})
		}
	}
	link.OriginalURL = longURL.OriginalURL

	err = s.db.EditLink(link)
	if err != nil {
		log.Printf("%v | %s", time.Now(), err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "something went wrong"})
	}
	return c.Status(fiber.StatusAccepted).JSON(linksResponse)
}

func (s *FiberServer) EnableDisbaleURLHandler(c *fiber.Ctx) error {
	sessionHeader := c.Get("Authorization")

	// ensure the session header is not empty and in the correct format
	if sessionHeader == "" || len(sessionHeader) < 8 || sessionHeader[:7] != "Bearer " {
		return c.Status(400).JSON(fiber.Map{"error": "invalid session header"})
	}
	// get the session id
	sessionId := sessionHeader[7:]

	_, err := s.GetSession("session:" + sessionId)
	if err != nil {
		log.Printf("%v | %s", time.Now(), err.Error())

		c.SendStatus(401)
		return c.JSON(fiber.Map{"message": "You are not logged in..."})
	}

	shortCode := c.Params("shortCode")

	val, err := strconv.ParseBool(c.Params("val"))
	if err != nil {
		c.SendStatus(400)
		return c.JSON(fiber.Map{"message": "bad rerquest"})
	}

	link, err := s.db.GetLink(shortCode)
	linksResponse := []types.LinkResponse{}
	if err != nil {
		if err == sql.ErrNoRows {
			c.SendStatus(fiber.StatusAccepted)
			return err
		} else {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "something went wrong"})
		}
	}

	link.IsEnabled = val
	err = s.db.EnableDisableLink(link)
	if err != nil {
		log.Printf("%v | %s", time.Now(), err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "something went wrong"})
	}
	return c.Status(fiber.StatusAccepted).JSON(linksResponse)
}
