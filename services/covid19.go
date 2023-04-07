package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/FaresAbuIram/COVID19-Statistics/entity"
	"github.com/FaresAbuIram/COVID19-Statistics/graph/model"
	"github.com/FaresAbuIram/COVID19-Statistics/logger"
)

type Covid19Service struct {
	SQLRepository    SQLRepository
	LoggerCollection logger.LoggerCollection
}

func NewCovid19Service(sqlRepository SQLRepository, loggerCollection logger.LoggerCollection) *Covid19Service {
	return &Covid19Service{
		SQLRepository:    sqlRepository,
		LoggerCollection: loggerCollection,
	}
}

func (c *Covid19Service) AddCountry(name string, userId int) (bool, error) {
	c.LoggerCollection.AddInfoLogger("services," + "covid19.go," + "AddCountry Func")

	count, err := c.SQLRepository.UsersCountById(userId)
	if err != nil {
		c.LoggerCollection.AddErrorLogger(err.Error())
		return false, err
	}

	if count == 0 {
		c.LoggerCollection.AddErrorLogger(fmt.Sprintf("user with id %d doesn't exist", userId))
		return false, fmt.Errorf("user with id %d doesn't exist", userId)
	}

	count, err = c.SQLRepository.CountriesCountByname(name)
	if err != nil {
		c.LoggerCollection.AddErrorLogger(err.Error())
		return false, err
	}
	var countryId int
	if count == 0 {
		countryId, err = c.SQLRepository.InsertCountry(name)
		if err != nil {
			c.LoggerCollection.AddErrorLogger(err.Error())
			return false, err
		}
		err = c.SQLRepository.InsertStatistic(countryId)
		if err != nil {
			c.LoggerCollection.AddErrorLogger(err.Error())
			return false, err
		}
	} else {
		countryId, err = c.SQLRepository.GetCountryIdByName(name)
		if err != nil {
			c.LoggerCollection.AddErrorLogger(err.Error())
			return false, err
		}
	}

	err = c.SQLRepository.InsertIntoUsersCountries(userId, countryId)
	if err != nil {
		c.LoggerCollection.AddErrorLogger(err.Error())
		return false, err
	}

	return true, nil
}

func (c *Covid19Service) GetCountries(userId int) ([]*model.Country, error) {
	countries, err := c.SQLRepository.GetAllCountriesByUserId(userId)
	if err != nil {
		c.LoggerCollection.AddErrorLogger(err.Error())
		return nil, err
	}
	return countries, nil
}

func (c *Covid19Service) PercentageOfDeathToConfirmed(userId int, countryName string) (float64, error) {
	percentage, err := c.SQLRepository.GetPercentageOfDeathToConfirmedByCountryName(userId, countryName)
	if err != nil {
		c.LoggerCollection.AddErrorLogger(err.Error())
		return 0.0, err
	}
	return percentage, nil
}

func (c *Covid19Service) GetTopThreeCountries(userId int, status string) ([]*model.Country, error) {
	countries, err := c.SQLRepository.GetTopThreeCountriesByUserIdAndType(userId, status)
	if err != nil {
		c.LoggerCollection.AddErrorLogger(err.Error())
		return nil, err
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
	countries, err := c.SQLRepository.GetAllCountries()
	if err != nil {
		return err
	}
	statistics, err := c.SQLRepository.GetAllStatistics()
	if err != nil {
		return err
	}
	newStatistics := c.fetchDataFromAPI(countries, statistics)
	c.SQLRepository.UpdateArrayOfStatistics(newStatistics)

	return nil
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

		var covidDataArray []entity.CovidData
		err = json.NewDecoder(resp.Body).Decode(&covidDataArray)
		if err != nil {
			fmt.Println(err)
			c.LoggerCollection.AddErrorLogger(err.Error())
			continue
		}
		if len(covidDataArray) == 0 {
			continue
		}

		statistics[index].Confirmed = int(covidDataArray[0].Confirmed)
		statistics[index].Deaths = int(covidDataArray[0].Deaths)
		statistics[index].Recovered = int(covidDataArray[0].Recovered)
	}

	return statistics
}
