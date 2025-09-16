package cmd

import (
	"github.com/LordVillain/Recommendation-service/internal/recommendation"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	mlURL := "http://localhost:8000" // Python ML сервис
    
	recSvc := recommendation.NewRecommendationService(mlURL)
	recommendation.NewRecommendationHandler(r, recSvc)

	r.Run(":8080")
}