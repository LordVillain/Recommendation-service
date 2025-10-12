package recommendation

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type RecommendationService struct {
	mlServiceURL string
}

func NewRecommendationService(mlURL string) *RecommendationService {
	return &RecommendationService{mlServiceURL: mlURL}
}

func (s *RecommendationService) GetRecommendations(req RecommendationRequest) (*RecommendationResponse, error) {
	resp, err := http.Get(fmt.Sprintf("%s/api/predict/%d/", s.mlServiceURL, req.UserID))
	if err != nil {
		return nil, fmt.Errorf("ML service unavailable: %v", err)
	}
	defer resp.Body.Close()

	// var recs RecommenderResponse

	var otvet RecommendationResponse

	if err := json.NewDecoder(resp.Body).Decode(&otvet); err != nil {
		return nil, fmt.Errorf("failed to decode ML response list of products: %v", err)
	}

	// TODO: grpc

	return &otvet, nil
}