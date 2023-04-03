package services

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/FaresAbuIram/COVID19-Statistics/logger"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	DB               *sql.DB
	LoggerCollection logger.LoggerCollection
}

func NewUserService(db *sql.DB, loggerCollection logger.LoggerCollection) *UserService {
	return &UserService{
		DB: db,
		LoggerCollection: loggerCollection,
	}
}

func (u *UserService) CreateNewUser(email, password string) (bool, error) {
	var count int
	err := u.DB.QueryRow("SELECT COUNT(*) FROM users WHERE email = $1", email).Scan(&count)
	if err != nil {
		return false, err
	}
	if count > 0 {
		return false, fmt.Errorf("user with email %s already exists", email)
	}

	// Hash the password with bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return false, err
	}

	// Insert the new user into the database
	_, err = u.DB.Exec("INSERT INTO users (email, password) VALUES ($1, $2)", email, hashedPassword)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (u *UserService) Login(email, password string) (string, error) {
	// Find the user with the given email address
	var id int64
	var hashedPassword []byte
	err := u.DB.QueryRow("SELECT id, password FROM users WHERE email = $1", email).Scan(&id, &hashedPassword)
	if err != nil {
		return "", fmt.Errorf("user with email %s not found", email)
	}

	// Check if the provided password matches the stored password
	if err := bcrypt.CompareHashAndPassword(hashedPassword, []byte(password)); err != nil {
		return "", fmt.Errorf("invalid password")
	}

	// Generate a new JWT token for the user
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": id,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	// Sign the token with the secret key
	tokenString, err := token.SignedString([]byte("secreatetoken"))
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	// // Insert the session token into the database
	// _, err = r.DB.Exec("INSERT INTO sessions (user_id, token) VALUES (?, ?)", id, tokenString)
	// if err != nil {
	// 	return "", err
	// }

	return tokenString, nil
}
