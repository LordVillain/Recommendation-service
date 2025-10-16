package recommendation

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/ShopOnGO/ShopOnGO/pkg/logger"
	pb "github.com/ShopOnGO/product-proto/pkg/product"
	"github.com/lib/pq"
	"github.com/shopspring/decimal"
	"google.golang.org/grpc"
)

type GRPCClients struct {
	ProductClient	pb.ProductServiceClient
}

type RecommendationService struct {
	mlServiceURL string
	Clients *GRPCClients
}

func NewRecommendationService(mlURL string) *RecommendationService {
	return &RecommendationService{
		mlServiceURL: mlURL,
		Clients: InitGRPCClients(),
	}
}

func InitGRPCClients() *GRPCClients {
	conn, err := grpc.Dial(fmt.Sprintf("%s:%s", os.Getenv("PRODUCT_SERVICE_HOST"), os.Getenv("PRODUCT_SERVICE_PORT")), grpc.WithInsecure())
	if err != nil {
		logger.Errorf("Ошибка подключения к gRPC серверу: %v", err)
	}

	logger.Info("gRPC connected")
	productClient := pb.NewProductServiceClient(conn)

	return &GRPCClients{
		ProductClient: productClient,
	}
}

func (s *RecommendationService) GetRecommendations(userID uint) (*RecommendationResponse, error) {
	ctx := context.Background()

	logger.Infof("Requesting products for user: %d", userID)
	resp, err := http.Get(fmt.Sprintf("%s/api/predict/%d/", s.mlServiceURL, userID))
	if err != nil {
		return nil, fmt.Errorf("ML service unavailable: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ML service returned status %d", resp.StatusCode)
	}

	var recs RecommenderResponse
	if err := json.NewDecoder(resp.Body).Decode(&recs); err != nil {
		return nil, fmt.Errorf("failed to decode ML response list of products: %v", err)
	}

	productIDsUint64 := make([]uint64, len(recs.ProductIDs))
	for i, id := range recs.ProductIDs {
		if id <= 0 {
			continue
		}
		productIDsUint64[i] = uint64(id)
	}

	logger.Infof("Product IDs: %v", productIDsUint64)

	// grpcResp, err := s.Clients.ProductClient.GetProductsByIDs(ctx, &pb.GetProductsByIDsRequest{
	// 	ProductIds: recs.ProductIDs,
	// })

	// добавить проверку что такие id продуктов есть

	grpcResp, err := s.Clients.ProductClient.GetProductsByIDs(ctx, &pb.GetProductsByIDsRequest{
    	ProductIds: productIDsUint64,
	})
	if err != nil {
		return nil, fmt.Errorf("error fetching products: %v", err)
	}

	logger.Info(grpcResp.Products)

	localProducts := make([]Product, len(grpcResp.Products))
	for i, p := range grpcResp.Products {
		rating := decimal.NewFromFloat(p.Rating)
		localProducts[i] = Product{
			Name:           p.Name,
			Description:    p.Description,
			Rating:         rating,
			ReviewCount:    uint(p.ReviewCount),
			RatingSum:      uint(p.RatingSum),
			QuestionCount:  uint(p.QuestionCount),
			IsActive:       p.IsActive,
			CategoryID:     uint(p.CategoryId),
			BrandID:        uint(p.BrandId),
			ImageURLs:      pq.StringArray(p.ImageUrls),
			VideoURLs:      pq.StringArray(p.VideoUrls),
		}
	}

	return &RecommendationResponse{
		UserID:   uint(userID),
		Products: localProducts,
	}, nil
}