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

func (s *RecommendationService) GetRecommendations(req RecommendationRequest) (*RecommendationResponse, error) {
	ctx := context.Background()

	resp, err := http.Get(fmt.Sprintf("%s/api/predict/%d/", s.mlServiceURL, req.UserID))
	if err != nil {
		return nil, fmt.Errorf("ML service unavailable: %v", err)
	}
	defer resp.Body.Close()

	var recs RecommenderResponse

	if err := json.NewDecoder(resp.Body).Decode(&recs); err != nil {
		return nil, fmt.Errorf("failed to decode ML response list of products: %v", err)
	}

	// grpcResp, err := s.Clients.ProductClient.GetProductsByIDs(ctx, &pb.GetProductsByIDsRequest{
	// 	ProductIds: recs.ProductIDs,
	// })
	grpcResp, err := s.Clients.ProductClient.GetProductsByIDs(ctx, &pb.GetProductsByIDsRequest{
    	ProductIds: []uint64{3, 4, 5},
	})
	if err != nil {
		return nil, fmt.Errorf("error fetching products: %v", err)
	}
	fmt.Println(grpcResp.Products)

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
		UserID:   req.UserID,
		Products: localProducts,
	}, nil
}