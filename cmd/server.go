// cmd/recommendation-service/main.go
package main

import (
	"fmt"

	"github.com/LordVillain/Recommendation-service/internal/recommendation"
	"github.com/LordVillain/Recommendation-service/pkg/middleware"
	"github.com/LordVillain/Recommendation-service/configs"
	"github.com/gin-gonic/gin"
)

func main() {
	config := configs.LoadConfig()

	router := gin.Default()
	fmt.Println("Starting recommendation service...")

	router.Use(middleware.GinAuthMiddleware(config))

	mlURL := "http://recommender-service:5000"
	recSvc := recommendation.NewRecommendationService(mlURL)
	recommendation.NewRecommendationHandler(router, recSvc)

	router.Run(":8086")
}