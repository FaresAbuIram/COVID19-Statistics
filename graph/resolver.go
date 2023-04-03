//go:generate go run github.com/99designs/gqlgen generate

package graph

import (
	"github.com/FaresAbuIram/COVID19-Statistics/services"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	UserService *services.UserService
	Covid19Service *services.Covid19Service
}
