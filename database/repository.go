package database

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/FaresAbuIram/COVID19-Statistics/entity"
	"github.com/FaresAbuIram/COVID19-Statistics/graph/model"
)

type SQLRepositoryInterface interface {
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

type SQLRepository struct {
	DB *sql.DB
}

func NewSQLRepository(db *sql.DB) *SQLRepository {
	return &SQLRepository{
		DB: db,
	}
}

func (sq *SQLRepository) UsersCountById(userId int) (int, error) {
	var count int
	err := sq.DB.QueryRow("SELECT COUNT(*) FROM users WHERE id = $1", userId).Scan(&count)
	return count, err
}

func (sq *SQLRepository) CountriesCountByname(name string) (int, error) {
	var count int
	err := sq.DB.QueryRow("SELECT COUNT(*) FROM countries WHERE name = $1", name).Scan(&count)
	return count, err
}

func (sq *SQLRepository) InsertCountry(name string) (int, error) {
	var countryId int
	err := sq.DB.QueryRow("INSERT INTO countries (name) VALUES ($1) RETURNING id", name).Scan(&countryId)
	return countryId, err
}

func (sq *SQLRepository) InsertStatistic(countryId int) error {
	_, err := sq.DB.Exec("INSERT INTO statistics (country_id) VALUES ($1)", countryId)
	return err
}

func (sq *SQLRepository) GetCountryIdByName(name string) (int, error) {
	var countryId int
	err := sq.DB.QueryRow("SELECT id from countries WHERE name = $1", name).Scan(&countryId)
	return countryId, err
}

func (sq *SQLRepository) InsertIntoUsersCountries(userId, countryId int) error {
	_, err := sq.DB.Exec("INSERT INTO users_countries (user_id, country_id) VALUES ($1, $2)", userId, countryId)
	return err
}

func (sq *SQLRepository) GetAllCountriesByUserId(userId int) ([]*model.Country, error) {
	rows, err := sq.DB.Query("SELECT name FROM countries WHERE id IN (SELECT country_id FROM users_countries WHERE user_id = $1);", userId)
	if err != nil {
		return nil, err
	}

	countries := make([]*model.Country, 0)
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		countries = append(countries, &model.Country{Name: name})
	}

	return countries, nil
}

func (sq *SQLRepository) GetPercentageOfDeathToConfirmedByCountryName(userId int, countryName string) (float64, error) {
	countryId, err := sq.GetCountryIdByName(countryName)
	if err != nil {
		return 0.0, err
	}
	query := `SELECT
					 (CAST(statistics.death AS FLOAT) / statistics.confirmed) * 100 AS death_to_confirmed_ratio 
			  FROM users_countries JOIN statistics ON users_countries.country_id = statistics.country_id 
			  WHERE users_countries.user_id = $1 AND users_countries.country_id = $2;
	`
	var percentage float64
	err = sq.DB.QueryRow(query, userId, countryId).Scan(&percentage)
	if err != nil {
		return 0.0, err
	}

	return percentage, nil
}

func (sq *SQLRepository) GetTopThreeCountriesByUserIdAndType(userId int, status string) ([]*model.Country, error) {
	query := `SELECT
					countries.name
				FROM
					users_countries
					JOIN countries ON users_countries.country_id = countries.id
					JOIN statistics ON users_countries.country_id = statistics.country_id
				WHERE
					users_countries.user_id = $1
				ORDER BY
					CASE
					WHEN $2 = 'confirmed' THEN statistics.confirmed
					WHEN $3 = 'death' THEN statistics.death
					END DESC
				LIMIT
					3  
				`

	rows, err := sq.DB.Query(query, userId, status, status)
	if err != nil {
		return nil, err
	}

	countries := make([]*model.Country, 0)
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		countries = append(countries, &model.Country{Name: name})
	}
	return countries, nil
}

func (sq *SQLRepository) GetAllCountries() (map[int]string, error) {
	// get all countries
	rows, err := sq.DB.Query("SELECT * from countries")
	if err != nil {
		return nil, err
	}

	countries := make(map[int]string)
	for rows.Next() {
		var country entity.Country
		if err := rows.Scan(&country.ID, &country.Name); err != nil {
			return nil, err
		}
		countries[country.ID] = country.Name
	}

	return countries, nil
}

func (sq *SQLRepository) GetAllStatistics() ([]entity.Statistics, error) {
	// get all statistics
	rows, err := sq.DB.Query("SELECT * from statistics")
	if err != nil {
		return []entity.Statistics{}, err
	}

	statistics := make([]entity.Statistics, 0)
	for rows.Next() {
		var statistic entity.Statistics
		if err := rows.Scan(&statistic.CountryId, &statistic.Confirmed, &statistic.Deaths, &statistic.Recovered, &statistic.LastUpdated); err != nil {
			return []entity.Statistics{}, err
		}
		statistics = append(statistics, statistic)
	}

	return statistics, nil
}

func (sq *SQLRepository) UpdateArrayOfStatistics(statistics []entity.Statistics) {
	// Update the database with the data
	for _, statistic := range statistics {
		_, err := sq.DB.Exec("UPDATE statistics SET confirmed = $1, death = $2, recovered = $3 WHERE country_id = $4", statistic.Confirmed, statistic.Deaths, statistic.Recovered, statistic.CountryId)
		if err != nil {
			continue
		}
	}
}

func (sq *SQLRepository) UsersCountByEmail(email string) (int, error) {
	var count int
	err := sq.DB.QueryRow("SELECT COUNT(*) FROM users WHERE email = $1", email).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (sq *SQLRepository) InsertNewUser(email string, password []byte) error {
	// Insert the new user into the database
	_, err := sq.DB.Exec("INSERT INTO users (email, password) VALUES ($1, $2)", email, password)
	return err
}

func (sq *SQLRepository) FindUserByEmail(email string) (int, []byte, error) {
	// Find the user with the given email address
	var id int
	var hashedPassword []byte
	err := sq.DB.QueryRow("SELECT id, password FROM users WHERE email = $1", email).Scan(&id, &hashedPassword)
	if err != nil {
		return 0, nil, errors.New(fmt.Sprintf("user with email %s not found", email))
	}

	return id, hashedPassword, nil
}
