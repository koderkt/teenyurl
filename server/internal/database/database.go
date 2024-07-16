package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/joho/godotenv/autoload"
	"github.com/koderkt/teenyurl/internal/types"
	_ "github.com/lib/pq"
)

// Service represents a service that interacts with a database.
type Service interface {
	Health() map[string]string
	Close() error
	CreateUser(*types.User) error
	Init() error
	GetUserByEmail(string) (*types.User, error)
	CreateShortURL(*types.Link) error
	GetLink(string) (*types.Link, error)
	GetLinks(int) (*[]types.Link, error)
	InsertAnalytics(*types.Clicks) error
	GetAnalystics(string) (*[]types.Clicks, error)
	GetNumberOfClicks(string) (int, error)
	// GetOriginalURL(string) (string, error)
}

type service struct {
	db *sqlx.DB
}

var (
	database = os.Getenv("DB_DATABASE")
	password = os.Getenv("DB_PASSWORD")
	username = os.Getenv("DB_USERNAME")
	// port       = os.Getenv("DB_PORT")
	// host       = os.Getenv("DB_HOST")
	// schema     = os.Getenv("DB_SCHEMA")
	dbInstance *service
)

func New() Service {
	// Reuse Connection
	if dbInstance != nil {
		return dbInstance
	}
	// connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable&search_path=%s", username, password, host, port, database, schema)
	// db, err := sql.Open("pgx", connStr)

	// sqlx
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", username, password, database)
	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	dbInstance = &service{
		db: db,
	}
	return dbInstance
}

func (s *service) Init() error {
	return s.createTables()
}

// Health checks the health of the database connection by pinging the database.
// It returns a map with keys indicating various health statistics.
func (s *service) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	stats := make(map[string]string)

	// Ping the database
	err := s.db.PingContext(ctx)
	if err != nil {
		stats["status"] = "down"
		stats["error"] = fmt.Sprintf("db down: %v", err)
		log.Fatalf(fmt.Sprintf("db down: %v", err)) // Log the error and terminate the program
		return stats
	}

	// Database is up, add more statistics
	stats["status"] = "up"
	stats["message"] = "It's healthy"

	// Get database stats (like open connections, in use, idle, etc.)
	dbStats := s.db.Stats()
	stats["open_connections"] = strconv.Itoa(dbStats.OpenConnections)
	stats["in_use"] = strconv.Itoa(dbStats.InUse)
	stats["idle"] = strconv.Itoa(dbStats.Idle)
	stats["wait_count"] = strconv.FormatInt(dbStats.WaitCount, 10)
	stats["wait_duration"] = dbStats.WaitDuration.String()
	stats["max_idle_closed"] = strconv.FormatInt(dbStats.MaxIdleClosed, 10)
	stats["max_lifetime_closed"] = strconv.FormatInt(dbStats.MaxLifetimeClosed, 10)

	// Evaluate stats to provide a health message
	if dbStats.OpenConnections > 40 { // Assuming 50 is the max for this example
		stats["message"] = "The database is experiencing heavy load."
	}

	if dbStats.WaitCount > 1000 {
		stats["message"] = "The database has a high number of wait events, indicating potential bottlenecks."
	}

	if dbStats.MaxIdleClosed > int64(dbStats.OpenConnections)/2 {
		stats["message"] = "Many idle connections are being closed, consider revising the connection pool settings."
	}

	if dbStats.MaxLifetimeClosed > int64(dbStats.OpenConnections)/2 {
		stats["message"] = "Many connections are being closed due to max lifetime, consider increasing max lifetime or revising the connection usage pattern."
	}

	return stats
}

func (s *service) Close() error {
	log.Printf("Disconnected from database: %s", database)
	return s.db.Close()
}

func (s *service) CreateUser(user *types.User) error {
	createUserQuery := `insert into users
	(user_name, email, encrypted_password, created_at)
	values ($1, $2, $3, $4)`

	userFromDb := &types.User{}

	getUserQuery := "select * from users where email = $1"
	err := s.db.Get(userFromDb, getUserQuery, user.Email)

	if err != sql.ErrNoRows {
		return errors.New("email/username already exists")
	}

	getUserQuery = "select * from users where user_name = $1"
	err = s.db.Get(userFromDb, getUserQuery, user.Email)

	if err != sql.ErrNoRows {
		return errors.New("email/username already exists")
	}
	_, err = s.db.Query(
		createUserQuery,
		user.UserName,
		user.Email,
		user.EncryptedPassword,
		user.CreatedAt,
	)

	if err != nil {
		return err
	}
	return nil
}

func (s *service) GetUserByEmail(email string) (*types.User, error) {
	userFromDb := &types.User{}
	getUserQuery := "select * from users where email = $1"
	err := s.db.Get(userFromDb, getUserQuery, email)

	if err == sql.ErrNoRows {
		return nil, errors.New("invalid email")
	}

	return userFromDb, nil
}

func (s *service) createTables() error {
	userTableQuery := `create table if not exists users (
		id serial primary key,
		user_name varchar(100),
		email varchar(100),
		encrypted_password varchar(100),
		created_at timestamp
	);`
	_, err := s.db.Exec(userTableQuery)
	// return err
	if err != nil {
		log.Fatalf("error while creating users table: %s", err.Error())

	}

	linkTableQuery := `CREATE TABLE if not exists urls (
    	id SERIAL PRIMARY KEY,
    original_url TEXT NOT NULL,
    short_url TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    user_id INT NOT NULL,
    is_enabled BOOLEAN NOT NULL DEFAULT TRUE
	);`
	_, err = s.db.Exec(linkTableQuery)
	// return err
	if err != nil {
		log.Fatalf("error while creating link table: %s", err.Error())
	}

	query := `CREATE TABLE IF NOT EXISTS clicks (
		id SERIAL PRIMARY KEY,
		short_code VARCHAR(6) NOT NULL,
		time_stamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		device_type VARCHAR(50),
		location VARCHAR(100)
		);`
	_, err = s.db.Exec(query)
	// return err
	if err != nil {
		log.Fatalf("error while creating clicks table: %s", err.Error())
	}

	return nil
}

func (s *service) CreateShortURL(link *types.Link) error {
	createLinkQuery := `insert into urls
	(original_url, short_url, user_id)
	values ($1, $2, $3)`

	_, err := s.db.Query(
		createLinkQuery,
		link.OriginalURL,
		link.ShortURL,
		link.UserId,
	)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) GetLink(shortURL string) (*types.Link, error) {
	record := &types.Link{}
	getLinkQuery := "select * from urls where short_url = $1"
	err := s.db.Get(record, getLinkQuery, shortURL)

	if err == sql.ErrNoRows {
		return nil, err
	}

	return record, nil
}

func (s *service) InsertAnalytics(analytics *types.Clicks) error {
	fmt.Println(*analytics)
	query := `INSERT INTO clicks (short_code, device_type, location)
	values ($1, $2, $3)`

	_, err := s.db.Query(
		query,
		analytics.ShortCode,
		analytics.DeviceType,
		analytics.Location,
	)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) GetAnalystics(shortCode string) (*[]types.Clicks, error) {
	record := &[]types.Clicks{}
	getClicksQuery := "select * from clicks where short_code = $1"
	err := s.db.Select(record, getClicksQuery, shortCode)

	if err != nil {
		return nil, err
	}

	return record, nil
}

func (s *service) GetLinks(userId int) (*[]types.Link, error) {
	var links []types.Link
	getClicksQuery := "SELECT * FROM urls WHERE user_id = $1"
	err := s.db.Select(&links, getClicksQuery, userId)

	if err != nil {
		return nil, err
	}

	return &links, nil
}

func (s *service) GetNumberOfClicks(shortURL string) (int, error) {
	query := `SELECT 
                COUNT(c.id) as click_count
              FROM 
                urls u
              LEFT JOIN 
                clicks c 
              ON 
                u.short_url = c.short_code
              WHERE 
                u.short_url = $1
              GROUP BY 
                u.short_url;`

    var clickCount int
    err := s.db.Get(&clickCount, query, shortURL)
    if err != nil {
        return 0, err
    }
	
    return clickCount, nil

}
