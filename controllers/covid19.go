package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/FaresAbuIram/COVID19-Statistics/entity"
	"github.com/FaresAbuIram/COVID19-Statistics/graph"
	"github.com/FaresAbuIram/COVID19-Statistics/logger"
	"github.com/FaresAbuIram/COVID19-Statistics/middleware"
	"github.com/gin-gonic/gin"
)

type Covid19Controller struct {
	Resolver *graph.Resolver
	Logger   logger.LoggerCollection
}

func NewCovid19Controller(resolver *graph.Resolver, logger logger.LoggerCollection) *Covid19Controller {
	return &Covid19Controller{
		Resolver: resolver,
		Logger:   logger,
	}
}

// Add new country
// @Summary      Add new country
// @Description  Add new country for a the user
// @Accept       json
// @Produce      json
// @Param		 Authorization	header		string	true	"Authentication header"
// @Param        body body entity.AddCountryRequest true "country name"
// @Success      200  {object}  entity.RegisterResponseSuccess
// @Failure      500  {object}	entity.UserResponseFailure
// @Router       /country [post]
func (cc *Covid19Controller) AddNewCountry(context *gin.Context) {
	cc.Logger.AddInfoLogger("controllers," + "covid19.go," + "AddNewCountry() Func")

	var userInput entity.AddCountryRequest
	if err := context.BindJSON(&userInput); err != nil {
		cc.Logger.AddErrorLogger(err.Error())
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	
	userId := middleware.GetUserID(context)
	fmt.Println(userId)
	mutation := `mutation {
		addCountry(input: {
			name: "%s",
			userId: %d
		}) 
	  }`
	query := fmt.Sprintf(mutation, userInput.Name, userId)
	fmt.Println(query)
	// Marshal the mutation variables to JSON
	bodyBody, _ := json.Marshal(map[string]string{
		"query": query,
	})

	// Create a new HTTP request to the GraphQL server
	resp, err := newQueryRequest(bodyBody)
	if err != nil {
		cc.Logger.AddErrorLogger(err.Error())
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	defer resp.Body.Close()

	var res struct {
		Data struct {
			AddCountry bool `json:"addCountry"`
		} `json:"data"`
	}

	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		cc.Logger.AddErrorLogger(err.Error())
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if res.Data.AddCountry == false {
		cc.Logger.AddErrorLogger("either country already exists or something wrong")
		context.JSON(http.StatusInternalServerError, gin.H{"error": "either country already exists or something wrong"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "country added successfully"})
}
