package services_test

import (
	"testing"

	"github.com/FaresAbuIram/COVID19-Statistics/graph/model"
	"github.com/FaresAbuIram/COVID19-Statistics/logger"
	"github.com/FaresAbuIram/COVID19-Statistics/services"
	SQLRepositoryInterface "github.com/FaresAbuIram/COVID19-Statistics/services/mocks"
)

func TestAddCountry(t *testing.T) {
	// prapare data
	sqlRepositoryInterface := new(SQLRepositoryInterface.SQLRepositoryInterface)
	logger := logger.NewLoggerCollection()
	covid19Service := services.NewCovid19Service(sqlRepositoryInterface, *logger)

	countryName := "Palestine"
	countryId := 1
	userId := 1

	sqlRepositoryInterface.On("UsersCountById", userId).Return(1, nil)
	sqlRepositoryInterface.On("CountriesCountByname",countryName).Return(1, nil)
	sqlRepositoryInterface.On("GetCountryIdByName",countryName).Return(1, nil)
	sqlRepositoryInterface.On("InsertIntoUsersCountries", userId, countryId).Return(nil)

	addCountry, err := covid19Service.AddCountry(countryName, userId)

	// Test cases
	if addCountry != true {
		t.Errorf("expected true adding country; got %v", addCountry)
	}

	// Test cases
	if err != nil {
		t.Errorf("expected nil error; got %v", err)
	}
}

func TestNegativeAddCountry(t *testing.T) {
	// prapare data
	sqlRepositoryInterface := new(SQLRepositoryInterface.SQLRepositoryInterface)
	logger := logger.NewLoggerCollection()
	covid19Service := services.NewCovid19Service(sqlRepositoryInterface, *logger)

	countryName := "Palestine"
	countryId := 1
	userId := 1

	sqlRepositoryInterface.On("UsersCountById", userId).Return(0, nil)
	sqlRepositoryInterface.On("CountriesCountByname",countryName).Return(1, nil)
	sqlRepositoryInterface.On("GetCountryIdByName",countryName).Return(1, nil)
	sqlRepositoryInterface.On("InsertIntoUsersCountries", userId, countryId).Return(nil)

	addCountry, err := covid19Service.AddCountry(countryName, userId)

	// Test cases
	if addCountry != false {
		t.Errorf("expected false adding country; got %v", addCountry)
	}

	// Test cases
	if err == nil {
		t.Errorf("expected got user with id 1 doesn't exist error; got %v", err)
	}
}

func TestGetCountries(t *testing.T) {
	// prapare data
	sqlRepositoryInterface := new(SQLRepositoryInterface.SQLRepositoryInterface)
	logger := logger.NewLoggerCollection()
	covid19Service := services.NewCovid19Service(sqlRepositoryInterface, *logger)

	userId := 1
	
	var countries []*model.Country
	countries = append(countries, &model.Country{Name: "Palestine"})
	countries = append(countries, &model.Country{Name: "Jordan"})

	sqlRepositoryInterface.On("GetAllCountriesByUserId", userId).Return(countries, nil)
	
	allCountries, err := covid19Service.GetCountries(userId)

	// Test cases
	if allCountries == nil {
		t.Errorf("expected array of countries; got %v", allCountries)
	}

	// Test cases
	if len(allCountries) != 2 {
		t.Errorf("expected 2 elements; got %v", len(allCountries))
	}

	
	// Test cases
	if allCountries[0].Name != "Palestine" {
		t.Errorf("expected Palestine; got %v", allCountries[0].Name)
	}

	// Test cases
	if err != nil {
		t.Errorf("expected nil error; got %v", err)
	}
}

func TestPercentageOfDeathToConfirmed(t *testing.T) {
	// prapare data
	sqlRepositoryInterface := new(SQLRepositoryInterface.SQLRepositoryInterface)
	logger := logger.NewLoggerCollection()
	covid19Service := services.NewCovid19Service(sqlRepositoryInterface, *logger)

	userId := 1
	companyName := "Palestine"

	sqlRepositoryInterface.On("GetPercentageOfDeathToConfirmedByCountryName", userId, companyName).Return(float64(10), nil)
	
	percentage, err := covid19Service.PercentageOfDeathToConfirmed(userId, companyName)

	// Test cases
	if percentage != float64(10) {
		t.Errorf("expected 10; got %v", percentage)
	}

	// Test cases
	if err != nil {
		t.Errorf("expected nil error; got %v", err)
	}
}

func TestGetTopThreeCountries(t *testing.T) {
	// prapare data
	sqlRepositoryInterface := new(SQLRepositoryInterface.SQLRepositoryInterface)
	logger := logger.NewLoggerCollection()
	covid19Service := services.NewCovid19Service(sqlRepositoryInterface, *logger)

	userId := 1
	status := "death"
	
	var countries []*model.Country
	countries = append(countries, &model.Country{Name: "Palestine"})
	countries = append(countries, &model.Country{Name: "Jordan"})
	countries = append(countries, &model.Country{Name: "Syria"})

	sqlRepositoryInterface.On("GetTopThreeCountriesByUserIdAndType", userId, status).Return(countries, nil)
	
	topThreeCountries, err := covid19Service.GetTopThreeCountries(userId, status)

	// Test cases
	if topThreeCountries == nil {
		t.Errorf("expected array of countries; got %v", topThreeCountries)
	}

	// Test cases
	if len(topThreeCountries) != 3 {
		t.Errorf("expected 3 elements; got %v", len(topThreeCountries))
	}

	
	// Test cases
	if topThreeCountries[2].Name != "Syria" {
		t.Errorf("expected Syria; got %v", topThreeCountries[2].Name)
	}

	// Test cases
	if err != nil {
		t.Errorf("expected nil error; got %v", err)
	}
}

