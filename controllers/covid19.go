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

	mutation := `mutation {
		addCountry(input: {
			name: "%s",
			userId: %d
		}) 
	  }`
	query := fmt.Sprintf(mutation, userInput.Name, userId)

	// Marshal the mutation variables to JSON
	queryBody, _ := json.Marshal(map[string]string{
		"query": query,
	})

	// Create a new HTTP request to the GraphQL server
	resp, err := newQueryRequest(queryBody)
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
		Errors []struct {
			Message string `json:"message"`
		} `json:"errors"`
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

	if len(res.Errors) != 0 {
		cc.Logger.AddErrorLogger(res.Errors[0].Message)
		context.JSON(http.StatusInternalServerError, gin.H{"error": res.Errors[0].Message})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "country added successfully"})
}

// Get all countries
// @Summary      Get all countries
// @Description  Get all countries subscribed by the user
// @Accept       json
// @Produce      json
// @Param		 Authorization	header		string	true	"Authentication header"
// @Success      200  {object}  []entity.CountryName
// @Failure      400  {object}	entity.UserResponseFailure
// @Failure      500  {object}	entity.UserResponseFailure
// @Router       /all-countries [get]
func (cc *Covid19Controller) GetCountries(context *gin.Context) {
	cc.Logger.AddInfoLogger("controllers," + "covid19.go," + "GetCountries() Func")

	userId := middleware.GetUserID(context)

	newQuery := `query {
					list(
						userId: %d
					){
						name
					}
				}`
	query := fmt.Sprintf(newQuery, userId)

	// Marshal the mutation variables to JSON
	queryBody, _ := json.Marshal(map[string]string{
		"query": query,
	})

	// Create a new HTTP request to the GraphQL server
	resp, err := newQueryRequest(queryBody)
	if err != nil {
		cc.Logger.AddErrorLogger(err.Error())
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	defer resp.Body.Close()

	var res struct {
		Data struct {
			List []entity.CountryName `json:"list"`
		} `json:"data"`
		Errors []struct {
			Message string `json:"message"`
		} `json:"errors"`
	}

	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		cc.Logger.AddErrorLogger(err.Error())
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(res.Errors) != 0 {
		cc.Logger.AddErrorLogger(res.Errors[0].Message)
		context.JSON(http.StatusInternalServerError, gin.H{"error": res.Errors[0].Message})
		return
	}

	context.JSON(http.StatusOK, gin.H{"countries": res.Data.List})
}

// Get the percentage of death cases to confirmed cases for a given country.
// @Summary      get the percentage of death cases to confirmed cases for a given country.
// @Description  get the percentage of death cases to confirmed cases for a given country.
// @Accept       json
// @Produce      json
// @Param		 Authorization	header		string	true	"Authentication header"
// @Param        name  path string true "country name"
// @Success      200  {object}  entity.Percentage
// @Failure      400  {object}	entity.UserResponseFailure
// @Failure      500  {object}	entity.UserResponseFailure
// @Router       /percentage-of-death-to-confirmed/{name} [get]
func (cc *Covid19Controller) PercentageOfDeathToConfirmed(context *gin.Context) {
	cc.Logger.AddInfoLogger("controllers," + "covid19.go," + "PercentageOfDeathToConfirmed() Func")
	name := context.Param("name")
	if name == "" {
		cc.Logger.AddErrorLogger("missing name")
		context.JSON(http.StatusBadRequest, gin.H{"error": "missing name"})
		return
	}
	userId := middleware.GetUserID(context)

	newQuery := `query {
					percentageeOfDeathToConfirmed(input : {
							userId:  %d
							name: "%s"
						}
					)
				}
	  `
	query := fmt.Sprintf(newQuery, userId, name)

	// Marshal the mutation variables to JSON
	queryBody, _ := json.Marshal(map[string]string{
		"query": query,
	})

	// Create a new HTTP request to the GraphQL server
	resp, err := newQueryRequest(queryBody)
	if err != nil {
		cc.Logger.AddErrorLogger(err.Error())
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	defer resp.Body.Close()

	var res struct {
		Data struct {
			PercentageeOfDeathToConfirmed float64 `json:"percentageeOfDeathToConfirmed"`
		} `json:"data"`
		Errors []struct {
			Message string `json:"message"`
		} `json:"errors"`
	}

	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		cc.Logger.AddErrorLogger(err.Error())
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if len(res.Errors) != 0 {
		cc.Logger.AddErrorLogger(res.Errors[0].Message)
		context.JSON(http.StatusInternalServerError, gin.H{"error": res.Errors[0].Message})
		return
	}
	context.JSON(http.StatusOK, gin.H{"countries": res.Data.PercentageeOfDeathToConfirmed})
}

// Get Top Three Countries based on the case type passed by the user (confirmed, death)
// @Summary     Get Top Three Countries based on the case type passed by the user (confirmed, death)
// @Description  get the top 3 countries (among the subscribed countries) by the total number of cases based on the case type passed by the user (confirmed, death).
// @Accept       json
// @Produce      json
// @Param		 Authorization	header		string	true	"Authentication header"
// @Param        type  path string true "(confirmed, death)"
// @Success      200  {object}  []entity.CountryName
// @Failure      400  {object}	entity.UserResponseFailure
// @Failure      500  {object}	entity.UserResponseFailure
// @Router       /top-three-countries/{type} [get]
func (cc *Covid19Controller) GetTopThreeCountries(context *gin.Context) {
	cc.Logger.AddInfoLogger("controllers," + "covid19.go," + "GetTopThreeCountries() Func")
	status := context.Param("type")
	if status == "" {
		cc.Logger.AddErrorLogger("missing type")
		context.JSON(http.StatusBadRequest, gin.H{"error": "missing type"})
		return
	}
	userId := middleware.GetUserID(context)

	newQuery := `query {
					getTopThreeCountries(input :{
							userId: %d
							type: "%s"
						}
						){
							name
					}
				}
	  `
	query := fmt.Sprintf(newQuery, userId, status)

	// Marshal the mutation variables to JSON
	queryBody, _ := json.Marshal(map[string]string{
		"query": query,
	})

	// Create a new HTTP request to the GraphQL server
	resp, err := newQueryRequest(queryBody)
	if err != nil {
		cc.Logger.AddErrorLogger(err.Error())
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	defer resp.Body.Close()

	var res struct {
		Data struct {
			GetTopThreeCountries []entity.CountryName `json:"getTopThreeCountries"`
		} `json:"data"`
		Errors []struct {
			Message string `json:"message"`
		} `json:"errors"`
	}

	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		cc.Logger.AddErrorLogger(err.Error())
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if len(res.Errors) != 0 {
		cc.Logger.AddErrorLogger(res.Errors[0].Message)
		context.JSON(http.StatusInternalServerError, gin.H{"error": res.Errors[0].Message})
		return
	}
	context.JSON(http.StatusOK, gin.H{"countries": res.Data.GetTopThreeCountries})
}
