package main

import (
	"fmt"

	"github.com/LordVillain/Recommendation-service/internal/recommendation"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	fmt.Println("Starting recommendation service...")

	mlURL := "http://recommender-service:5000" // Python ML сервис
    
	recSvc := recommendation.NewRecommendationService(mlURL)
	recommendation.NewRecommendationHandler(router, recSvc)

	router.Run(":8086")
}