package main

import (
	"restaurant_golang/controllers"

	"github.com/gin-gonic/gin"

	_ "restaurant_golang/docs"

	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

// @title Swagger Example API
// @version 1.0
// @license.name CalvinR

func main() {
	router := gin.Default()
	v1 := router.Group("/apis")
	{
		organization := v1.Group("/organization")
		{
			organization.PUT("/:organizationId", controllers.UpdateOrganization)
			organization.GET("/:userId", controllers.GetOrganzation)
		}

		businessType := v1.Group("/businessType")
		{
			businessType.GET("/", controllers.GetBusinessType)
			businessType.GET("/subType/:id", controllers.GetBusinessSubType)
		}

		auth := v1.Group("/auth")
		{
			auth.POST("/registerUser", controllers.RegisterUser)
			auth.GET("/verifyEmail", controllers.VerifyEmail)
		}

		store := v1.Group("/store")
		{
			store.POST("/", controllers.CreateStore)
			store.GET("/:organizationId", controllers.GetStoreList)
		}

		notification := v1.Group("/notification")
		{
			notification.POST("/", controllers.SendMessaging)
		}
	}
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.Run(":8080")
}
