package services_test

import (
	"testing"

	"github.com/FaresAbuIram/COVID19-Statistics/logger"
	"github.com/FaresAbuIram/COVID19-Statistics/services"
	SQLRepositoryInterface "github.com/FaresAbuIram/COVID19-Statistics/services/mocks"
	"github.com/stretchr/testify/mock"
)

func TestCreateNewUser(t *testing.T) {
	// prapare data
	sqlRepositoryInterface := new(SQLRepositoryInterface.SQLRepositoryInterface)
	logger := logger.NewLoggerCollection()
	userService := services.NewUserService(sqlRepositoryInterface, *logger)

	fakeEmail := "test@test.com"
	sqlRepositoryInterface.On("UsersCountByEmail", fakeEmail).Return(0, nil)
	sqlRepositoryInterface.On("InsertNewUser", fakeEmail, mock.AnythingOfType("[]uint8")).Return(nil)

	register, err := userService.CreateNewUser(fakeEmail, "test")

	// Test cases
	if register != true {
		t.Errorf("expected true register; got %v", register)
	}

	// Test cases
	if err != nil {
		t.Errorf("expected nil error; got %v", err)
	}
}

func TestNegativeCreateNewUser(t *testing.T) {
	// prapare data
	sqlRepositoryInterface := new(SQLRepositoryInterface.SQLRepositoryInterface)
	logger := logger.NewLoggerCollection()
	userService := services.NewUserService(sqlRepositoryInterface, *logger)

	fakeEmail := "test@test.com"
	sqlRepositoryInterface.On("UsersCountByEmail", fakeEmail).Return(1, nil)
	sqlRepositoryInterface.On("InsertNewUser", fakeEmail, mock.AnythingOfType("[]uint8")).Return(nil)

	register, err := userService.CreateNewUser(fakeEmail, "test")

	// Test cases
	if register != false {
		t.Errorf("expected true register; got %v", register)
	}

	// Test cases
	if err == nil {
		t.Errorf("expected got user with email test@test.com already exists error; got %v", err)
	}
}

func TestNegativeLogin(t *testing.T) {
	// prapare data
	sqlRepositoryInterface := new(SQLRepositoryInterface.SQLRepositoryInterface)
	logger := logger.NewLoggerCollection()
	userService := services.NewUserService(sqlRepositoryInterface, *logger)

	fakeEmail := "test@test.com"
	fakePass := []byte("$2a$10$JEUwvw/FW8u.JnsW.v2YeOj6rQIN67wbom7cn578ydYLUjnO8RM5m")
	sqlRepositoryInterface.On("FindUserByEmail", fakeEmail).Return(1, fakePass, nil)

	token, err := userService.Login(fakeEmail, "test1")

	// Test cases
	if token != "" {
		t.Errorf("expected not empty token; got %v", token)
	}

	// Test cases
	if err == nil {
		t.Errorf("expected got invalid passwor error; got %v", err)
	}
}

func TestLogin(t *testing.T) {
	// prapare data
	sqlRepositoryInterface := new(SQLRepositoryInterface.SQLRepositoryInterface)
	logger := logger.NewLoggerCollection()
	userService := services.NewUserService(sqlRepositoryInterface, *logger)

	fakeEmail := "test@test.com"
	fakePass := []byte("$2a$10$JEUwvw/FW8u.JnsW.v2YeOj6rQIN67wbom7cn578ydYLUjnO8RM5m")
	sqlRepositoryInterface.On("FindUserByEmail", fakeEmail).Return(1, fakePass, nil)

	token, err := userService.Login(fakeEmail, "test")

	// Test cases
	if token == "" {
		t.Errorf("expected not empty token; got %v", token)
	}

	// Test cases
	if err != nil {
		t.Errorf("expected nil error; got %v", err)
	}
}