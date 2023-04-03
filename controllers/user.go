package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/FaresAbuIram/COVID19-Statistics/graph"
	"github.com/FaresAbuIram/COVID19-Statistics/graph/model"
	"github.com/FaresAbuIram/COVID19-Statistics/logger"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	Resolver *graph.Resolver
	Logger   logger.LoggerCollection
}

func NewUserController(resolver *graph.Resolver, logger logger.LoggerCollection) *UserController {
	return &UserController{
		Resolver: resolver,
		Logger:   logger,
	}
}

func (uc *UserController) Query(context *gin.Context) {
	h := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: uc.Resolver}))

	h.ServeHTTP(context.Writer, context.Request)
}

func newQueryRequest(queryBody []byte) (*http.Response, error) {
	request, err := http.NewRequest("POST", "http://localhost:8080/query", bytes.NewBuffer(queryBody))
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")
	// Send the HTTP request and read the response
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	
	return resp, nil
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
	uc.Logger.AddInfoLogger("controllers," + "user.go," + "Register() Func")

	var userInput model.RegisterInput
	if err := context.BindJSON(&userInput); err != nil {
		uc.Logger.AddErrorLogger(err.Error())
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	mutation := `mutation {
				register(input: {
					email: "%s",
					password: "%s"
				}) 
	  		}`
	query := fmt.Sprintf(mutation, userInput.Email, userInput.Password)
	
	// Marshal the mutation variables to JSON
	bodyBody, _ := json.Marshal(map[string]string{
		"query": query,
	})

	// Create a new HTTP request to the GraphQL server
	resp, err := newQueryRequest(bodyBody)
	if err != nil {
		uc.Logger.AddErrorLogger(err.Error())
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	defer resp.Body.Close()

	var res struct {
		Data struct {
			Register bool `json:"register"`
		} `json:"data"`
		
	}

	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		uc.Logger.AddErrorLogger(err.Error())
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if res.Data.Register == false {
		uc.Logger.AddErrorLogger(fmt.Sprintf("user with email %s already exists", userInput.Email))
		context.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("user with email %s already exists", userInput.Email)})
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
	uc.Logger.AddInfoLogger("controllers," + "user.go," + "Login() Func")

	var userInput model.LoginInput
	if err := context.BindJSON(&userInput); err != nil {
		uc.Logger.AddErrorLogger(err.Error())
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	mutation := `mutation {
				login(input: {
					email: "%s",
					password: "%s"
				}) 
	  		}`
	query := fmt.Sprintf(mutation, userInput.Email, userInput.Password)

	// Marshal the mutation variables to JSON
	bodyBody, _ := json.Marshal(map[string]string{
		"query": query,
	})

	// Create a new HTTP request to the GraphQL server
	resp, err := newQueryRequest(bodyBody)
	if err != nil {
		uc.Logger.AddErrorLogger(err.Error())
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	defer resp.Body.Close()

	var res struct {
		Data struct {
			Login string `json:"login"`
		} `json:"data"`
	}

	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		uc.Logger.AddErrorLogger(err.Error())
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if res.Data.Login == "" {
		uc.Logger.AddErrorLogger("user doesn't exist")
		context.JSON(http.StatusInternalServerError, gin.H{"error": "user doesn't exist"})
		return
	}
	

	context.JSON(http.StatusOK, gin.H{"token": res.Data.Login})
}
