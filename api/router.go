package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

// SetupRouter determine all the routes for this service
func SetupRouter() *gin.Engine {
	// setting Gin mode before running
	gin.SetMode(viper.GetString("GIN_MODE"))

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	//router.Use(middlewares.AuthMiddleware())
	validatorsRegistration()

	// healthCheck router
	health := new(HealthController)
	router.GET("/health", health.HealthCheck)

	v1 := router.Group("v1")
	{
		userGroup := v1.Group("user")
		{
			userController := NewUserController()
			//TODO : mettre les bonnes méthodes en face
			//userGroup.GET("/", userController.GetByID)
			userGroup.GET("/:id", userController.GetByID)
			userGroup.POST("/", userController.Create)
			userGroup.PUT("/:id", userController.CompleteUpdate)
			userGroup.PATCH("/:id", userController.PartialUpdate)
			userGroup.DELETE("/:id", userController.Delete)
		}
	}

	// Swagger endpoint
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Returning 404 when calling an unmapped uri
	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, ApiErrorResponse{
			Error: "Resource not found",
		})
	})

	return router
}
