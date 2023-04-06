package services

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/FaresAbuIram/COVID19-Statistics/entity"
	"github.com/FaresAbuIram/COVID19-Statistics/graph/model"
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
		c.LoggerCollection.AddErrorLogger(fmt.Sprintf("user with id %d doesn't exist", userId))
		return false, fmt.Errorf("user with id %d doesn't exist", userId)
	}

	err = c.DB.QueryRow("SELECT COUNT(*) FROM countries WHERE name = $1", name).Scan(&count)
	if err != nil {
		return false, err
	}
	var countryId int
	if count == 0 {

		err := c.DB.QueryRow("INSERT INTO countries (name) VALUES ($1) RETURNING id", name).Scan(&countryId)
		if err != nil {
			c.LoggerCollection.AddErrorLogger(err.Error())
			return false, err
		}

		_, err = c.DB.Exec("INSERT INTO statistics (country_id) VALUES ($1)", countryId)
		if err != nil {
			c.LoggerCollection.AddErrorLogger(err.Error())
			return false, err
		}
	} else {
		err := c.DB.QueryRow("SELECT id from countries WHERE name = $1", name).Scan(&countryId)
		if err != nil {
			c.LoggerCollection.AddErrorLogger(err.Error())
			return false, err
		}
	}

	_, err = c.DB.Exec("INSERT INTO users_countries (user_id, country_id) VALUES ($1, $2)", userId, countryId)
	if err != nil {
		fmt.Println(userId, countryId)
		c.LoggerCollection.AddErrorLogger(err.Error())
		return false, err
	}

	return true, nil
}

func (c *Covid19Service) GetCountries(userId int) ([]*model.Country, error) {
	rows, err := c.DB.Query("SELECT name FROM countries WHERE id IN (SELECT country_id FROM users_countries WHERE user_id = $1);", userId)
	if err != nil {
		c.LoggerCollection.AddErrorLogger(err.Error())
		return nil, err
	}

	countries := make([]*model.Country, 0)
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			c.LoggerCollection.AddErrorLogger(err.Error())
			return nil, err
		}
		countries = append(countries, &model.Country{Name: name})
	}
	return countries, nil
}

func (c *Covid19Service) PercentageOfDeathToConfirmed(userId int, countryName string) (float64, error) {
	var countryId int
	err := c.DB.QueryRow("SELECT id FROM countries WHERE name = $1", countryName).Scan(&countryId)
	if err != nil {
		fmt.Println(err)
		return 0.0, err
	}
	query := `SELECT
					 (CAST(statistics.death AS FLOAT) / statistics.confirmed) * 100 AS death_to_confirmed_ratio 
			  FROM users_countries JOIN statistics ON users_countries.country_id = statistics.country_id 
			  WHERE users_countries.user_id = $1 AND users_countries.country_id = $2;
	`
	var percentage float64
	err = c.DB.QueryRow(query, userId, countryId).Scan(&percentage)
	if err != nil {
		c.LoggerCollection.AddErrorLogger(err.Error())
		return 0.0, err
	}

	return percentage, nil
}

func (c *Covid19Service) GetTopThreeCountries(userId int, status string) ([]*model.Country, error) {
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

	rows, err := c.DB.Query(query, userId, status, status)
	if err != nil {
		c.LoggerCollection.AddErrorLogger(err.Error())
		return nil, err
	}

	countries := make([]*model.Country, 0)
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			c.LoggerCollection.AddErrorLogger(err.Error())
			return nil, err
		}
		countries = append(countries, &model.Country{Name: name})
	}
	return countries, nil
}

func (c *Covid19Service) GetDailyTotals() {
	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			c.fetchAndUpdateData()
		}
	}
}

func (c *Covid19Service) fetchAndUpdateData() error {
	countries, err := c.getCountriesFromDatabase()
	if err != nil {
		return err
	}
	statistics, err := c.getStatisticsFromDatabase()
	if err != nil {
		return err
	}
	newStatistics := c.fetchDataFromAPI(countries, statistics)
	c.updateStatistics(newStatistics)

	return nil
}

func (c *Covid19Service) getCountriesFromDatabase() (map[int]string, error) {
	// get all countries
	rows, err := c.DB.Query("SELECT * from countries")
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

func (c *Covid19Service) getStatisticsFromDatabase() ([]entity.Statistics, error) {
	// get all statistics
	rows, err := c.DB.Query("SELECT * from statistics")
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

func (c *Covid19Service) fetchDataFromAPI(countries map[int]string, statistics []entity.Statistics) []entity.Statistics {
	for index, statistic := range statistics {
		fromDate := statistic.LastUpdated.AddDate(-3, 0, 0)
		toDate := time.Date(fromDate.Year(), fromDate.Month(), fromDate.Day(), 23, 59, 59, 0, fromDate.Location())

		resp, err := http.Get(fmt.Sprintf("https://api.covid19api.com/total/country/%s?from=%s&to=%s", countries[statistic.CountryId], fromDate.UTC().Format("2006-01-02T00:00:00Z"), toDate.UTC().Format("2006-01-02T15:04:05Z")))
		if err != nil {
			c.LoggerCollection.AddErrorLogger(err.Error())
			continue
		}
		defer resp.Body.Close()

		var covidDataArray []interface{}
		err = json.NewDecoder(resp.Body).Decode(&covidDataArray)
		if err != nil {
			c.LoggerCollection.AddErrorLogger(err.Error())
			continue
		}

		latestData := covidDataArray[len(covidDataArray)-1].(map[string]interface{})
		statistics[index].Confirmed = int(latestData["Confirmed"].(float64))
		statistics[index].Deaths = int(latestData["Deaths"].(float64))
		statistics[index].Recovered = int(latestData["Recovered"].(float64))
	}

	return statistics
}

func (c *Covid19Service) updateStatistics(statistics []entity.Statistics) {
	// Update the database with the data
	for _, statistic := range statistics {
		_, err := c.DB.Exec("UPDATE statistics SET confirmed = $1, death = $2, recovered = $3 WHERE country_id = $4", statistic.Confirmed, statistic.Deaths, statistic.Recovered, statistic.CountryId)
		if err != nil {
			c.LoggerCollection.AddErrorLogger(err.Error())
			continue
		}
	}
}
