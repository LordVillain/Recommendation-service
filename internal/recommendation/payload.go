package recommendation


type RecommendationRequest struct {
	UserID uint `json:"user_id"`
}

type RecommenderResponse struct {
	ProductIDs []uint64 `json:"product_ids"`
}

type RecommendationResponse struct {
	UserID   uint       `json:"user_id"`
	Products []Product  `json:"products"`
}
