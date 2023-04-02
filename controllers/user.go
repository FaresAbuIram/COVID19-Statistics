package controllers

import (
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/FaresAbuIram/COVID19-Statistics/graph"
	"github.com/FaresAbuIram/COVID19-Statistics/graph/model"
	"github.com/FaresAbuIram/COVID19-Statistics/services"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	UserService *services.UserService
}

func NewUserController(userService *services.UserService) *UserController {
	return &UserController{
		UserService: userService,
	}
}

func (uc *UserController) Query(context *gin.Context) {
	h := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{
		UserService: uc.UserService,
	}}))

	h.ServeHTTP(context.Writer, context.Request)
}

// Create New User
// @Summary      Create New User
// @Description  Create New User with email and password
// @Accept       json
// @Produce      json
// @Param        body  body model.RegisterInput true "user's inforamtion"
// @Success      200  {object}  entity.RegisterResponseSuccess
// @Failure      500  {object}	entity.RegisterResponseFailure
// @Router       /register [post]
func (uc *UserController) Register(context *gin.Context) {
	var userInput model.RegisterInput
	if err := context.BindJSON(&userInput); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	resolver := &graph.Resolver{UserService: uc.UserService}
	_, err := resolver.UserService.CreateNewUser(userInput.Email, userInput.Password)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
}

func (uc *UserController) Login(context *gin.Context) {
	var userInput model.RegisterInput
	if err := context.BindJSON(&userInput); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	resolver := &graph.Resolver{UserService: uc.UserService}
	token, err := resolver.UserService.Login(userInput.Email, userInput.Password)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"token": token})
}
