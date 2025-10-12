package recommendation

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type RecommendationHandler struct {
	recommendationSvc *RecommendationService
}

func NewRecommendationHandler(router *gin.Engine, svc *RecommendationService) *RecommendationHandler {
	handler := &RecommendationHandler{recommendationSvc: svc}

	recGroup := router.Group("/recommendation-service/recommendations")
	{
		recGroup.POST("", handler.getRecommendations)
	}

	return handler
}

// getRecommendations godoc
// @Summary Получить рекомендации товаров для пользователя
// @Description Возвращает список рекомендованных товаров на основе истории покупок
// @Tags Рекомендации
// @Param request body models.RecommendationRequest true "Данные пользователя"
// @Success 200 {object} models.RecommendationResponse
// @Failure 400 {object} gin.H "Некорректные данные"
// @Failure 500 {object} gin.H "Ошибка сервиса"
// @Router /recommendations-service/recommendations [post]
func (h *RecommendationHandler) getRecommendations(c *gin.Context) {
	var req RecommendationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректные данные"})
		return
	}

	recs, err := h.recommendationSvc.GetRecommendations(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, recs)
}
