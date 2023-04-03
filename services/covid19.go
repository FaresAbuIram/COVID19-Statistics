package services

import (
	"database/sql"
	"fmt"

	"github.com/FaresAbuIram/COVID19-Statistics/logger"
)

type Covid19Service struct {
	DB               *sql.DB
	LoggerCollection logger.LoggerCollection
}

func NewCovid19Service(db *sql.DB, loggerCollection logger.LoggerCollection) *Covid19Service {
	return &Covid19Service{
		DB:               db,
		LoggerCollection: loggerCollection,
	}
}

func (c *Covid19Service) AddCountry(name string, userId int) (bool, error) {
	c.LoggerCollection.AddInfoLogger("services," + "covid19.go," + "AddCountry Func")

	var count int
	err := c.DB.QueryRow("SELECT COUNT(*) FROM users WHERE id = $1", userId).Scan(&count)
	if err != nil {
		return false, err
	}
	if count == 0 {
		c.LoggerCollection.AddErrorLogger(fmt.Sprintf("user with id %s doesn't exist", userId))
		return false, fmt.Errorf("user with id %s doesn't exist", userId)
	}

	// Insert the new user into the database
	_, err = c.DB.Exec("INSERT INTO country (user_id, name) VALUES ($1, $2)", userId, name)
	if err != nil {
		c.LoggerCollection.AddErrorLogger(err.Error())
		return false, err
	}

	return true, nil
}
