package recommendation

import (
	"bytes"
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
	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %v", err)
	}

	resp, err := http.Post(s.mlServiceURL+"/predict", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("ML service unavailable: %v", err)
	}
	defer resp.Body.Close()

	var recs RecommendationResponse
	if err := json.NewDecoder(resp.Body).Decode(&recs); err != nil {
		return nil, fmt.Errorf("failed to decode ML response: %v", err)
	}

	return &recs, nil
}