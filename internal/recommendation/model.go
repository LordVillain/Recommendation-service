package recommendation

type RecommendationRequest struct {
	UserID  uint   `json:"user_id"`
	History []string `json:"history"`
}

type RecommendationResponse struct {
	Recommendations []string `json:"recommendations"`
}

