package routes

import (
	"log"

	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/FaresAbuIram/COVID19-Statistics/controllers"
	"github.com/FaresAbuIram/COVID19-Statistics/database"
	docs "github.com/FaresAbuIram/COVID19-Statistics/docs"
	"github.com/FaresAbuIram/COVID19-Statistics/graph"
	"github.com/FaresAbuIram/COVID19-Statistics/logger"
	"github.com/FaresAbuIram/COVID19-Statistics/middleware"
	"github.com/FaresAbuIram/COVID19-Statistics/services"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Setup(router *gin.Engine) {
	docs.SwaggerInfo.Title = "Swagger Example API"
	docs.SwaggerInfo.Description = "This is a sample server Petstore server."
	docs.SwaggerInfo.Version = "2.0"
	docs.SwaggerInfo.Host = "localhost:8080"
	docs.SwaggerInfo.Schemes = []string{"http"}

	db, err := database.Connect()
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
		return
	}

	logger := logger.NewLoggerCollection()
	userService := services.NewUserService(db, *logger)
	covid19Service := services.NewCovid19Service(db, *logger)
	resolver := &graph.Resolver{UserService: userService, Covid19Service: covid19Service}
	userController := controllers.NewUserController(resolver, *logger)
	covid19Controller := controllers.NewCovid19Controller(resolver, *logger)

	go covid19Service.GetDailyTotals()
	router.Use(static.Serve("/", static.LocalFile("./website/dist", true)))
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.POST("/query", userController.Query)
	router.GET("/", gin.WrapH(playground.Handler("GraphQL playground", "/query")))
	router.POST("/register", userController.Register)
	router.POST("/login", userController.Login)
	router.POST("/country", middleware.AuthMiddleware(), covid19Controller.AddNewCountry)
	router.GET("/all-countries", middleware.AuthMiddleware(), covid19Controller.GetCountries)
	router.GET("/percentage-of-death-to-confirmed/:name", middleware.AuthMiddleware(), covid19Controller.PercentageOfDeathToConfirmed)
	router.GET("/top-three-countries/:type", middleware.AuthMiddleware(), covid19Controller.GetTopThreeCountries)

}
