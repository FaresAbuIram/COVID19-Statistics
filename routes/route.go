package routes

import (
	"log"

	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/FaresAbuIram/COVID19-Statistics/controllers"
	"github.com/FaresAbuIram/COVID19-Statistics/database"
	"github.com/FaresAbuIram/COVID19-Statistics/services"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/swag/example/basic/docs"
)

func Setup(router *gin.Engine) {
	docs.SwaggerInfo.Title = "Swagger Example API"
	docs.SwaggerInfo.Description = "This is a sample server Petstore server."
	docs.SwaggerInfo.Version = "2.0"
	docs.SwaggerInfo.Host = "localhost:8081"
	docs.SwaggerInfo.Schemes = []string{"http"}

	db, err := database.Connect()
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
		return
	}

	userService := services.NewUserService(db)
	userController := controllers.NewUserController(userService)

	// resourcedatabase := database.NewResourceService()
	// resourceService := controllers.NewResourceService(resourcedatabase)

	router.Use(static.Serve("/", static.LocalFile("./website/dist", true)))
	// Serve the GraphQL API at /graphql
	router.POST("/query", userController.Query)
	router.GET("/", gin.WrapH(playground.Handler("GraphQL playground", "/query")))
	router.POST("/register",userController.Register)
	router.POST("/login",userController.Login)
	// router.POST("/search", resourceService.Search)
	// router.GET("/resource/:id", resourceService.GetResource)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	
}
