package services

import (
	"fmt"
	"os"
	"time"

	"github.com/FaresAbuIram/COVID19-Statistics/entity"
	"github.com/FaresAbuIram/COVID19-Statistics/graph/model"
	"github.com/FaresAbuIram/COVID19-Statistics/logger"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

type SQLRepository interface {
	UsersCountById(userId int) (int, error)
	CountriesCountByname(name string) (int, error)
	InsertCountry(name string) (int, error)
	InsertStatistic(countryId int) error
	GetCountryIdByName(name string) (int, error)
	InsertIntoUsersCountries(userId, countryId int) error
	GetAllCountriesByUserId(userId int) ([]*model.Country, error)
	GetPercentageOfDeathToConfirmedByCountryName(userId int, countryName string) (float64, error)
	GetTopThreeCountriesByUserIdAndType(userId int, status string) ([]*model.Country, error)
	GetAllCountries() (map[int]string, error)
	GetAllStatistics() ([]entity.Statistics, error)
	UpdateArrayOfStatistics(statistics []entity.Statistics)
	UsersCountByEmail(email string) (int, error)
	InsertNewUser(email string, password []byte) error
	FindUserByEmail(email string) (int, []byte, error)
}
type UserService struct {
	SQLRepository    SQLRepository
	LoggerCollection logger.LoggerCollection
}

func NewUserService(sqlRepository SQLRepository, loggerCollection logger.LoggerCollection) *UserService {
	return &UserService{
		SQLRepository:    sqlRepository,
		LoggerCollection: loggerCollection,
	}
}

func (u *UserService) CreateNewUser(email, password string) (bool, error) {
	count, err := u.SQLRepository.UsersCountByEmail(email)
	if err != nil {
		u.LoggerCollection.AddErrorLogger(err.Error())
		return false, err
	}

	if count > 0 {
		u.LoggerCollection.AddErrorLogger(fmt.Errorf("user with email %s already exists", email).Error())
		return false, fmt.Errorf("user with email %s already exists", email)
	}

	// Hash the password with bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		u.LoggerCollection.AddErrorLogger(err.Error())
		return false, err
	}

	// Insert the new user into the database
	err = u.SQLRepository.InsertNewUser(email, hashedPassword)
	if err != nil {
		u.LoggerCollection.AddErrorLogger(err.Error())
		return false, err
	}

	return true, nil
}

func (u *UserService) Login(email, password string) (string, error) {
	id, hashedPassword, err := u.SQLRepository.FindUserByEmail(email)
	if err != nil {
		u.LoggerCollection.AddErrorLogger(err.Error())
		return "", fmt.Errorf("user with email %s not found", email)
	}

	// Check if the provided password matches the stored password
	if err := bcrypt.CompareHashAndPassword(hashedPassword, []byte(password)); err != nil {
		u.LoggerCollection.AddErrorLogger(err.Error())
		return "", fmt.Errorf("invalid password")
	}

	// Generate a new JWT token for the user
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": id,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	// Sign the token with the secret key
	tokenString, err := token.SignedString([]byte(os.Getenv("TOKEN_SECRET")))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
