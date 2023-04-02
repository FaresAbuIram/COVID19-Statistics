package controllers

import (
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/FaresAbuIram/COVID19-Statistics/graph"
	"github.com/FaresAbuIram/COVID19-Statistics/graph/model"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	Resolver *graph.Resolver
}

func NewUserController(resolver *graph.Resolver) *UserController {
	return &UserController{
		Resolver: resolver,
	}
}

func (uc *UserController) Query(context *gin.Context) {
	h := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: uc.Resolver}))

	h.ServeHTTP(context.Writer, context.Request)
}

// Create New User
// @Summary      Create New User
// @Description  Create New User with email and password
// @Accept       json
// @Produce      json
// @Param        body body model.RegisterInput true "email and password"
// @Success      200  {object}  entity.RegisterResponseSuccess
// @Failure      500  {object}	entity.UserResponseFailure
// @Router       /register [post]
func (uc *UserController) Register(context *gin.Context) {
	var userInput model.RegisterInput
	if err := context.BindJSON(&userInput); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	_, err := uc.Resolver.UserService.CreateNewUser(userInput.Email, userInput.Password)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
}

// Login
// @Summary    Login
// @Description  Login User with email and password
// @Accept       json
// @Produce      json
// @Param        body body model.LoginInput true "email and password"
// @Success      200  {object}  entity.LoginResponseSuccess
// @Failure      500  {object}	entity.UserResponseFailure
// @Router       /login [post]
func (uc *UserController) Login(context *gin.Context) {
	var userInput model.LoginInput
	if err := context.BindJSON(&userInput); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	token, err := uc.Resolver.UserService.Login(userInput.Email, userInput.Password)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"token": token})
}
