package stockms

import (
	"cart/internal/domain"
	"cart/internal/usecase/carts"
	pb "cart/pkg/api/stocks"
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	connTimeOut     = 10 * time.Second
	grpcCallTimeOut = 5 * time.Second
)

type grpcStockService struct {
	client pb.StocksServiceClient
	conn   *grpc.ClientConn
}

var _ carts.StockService = (*grpcStockService)(nil)

func NewGRPCStockService(address string) (*grpcStockService, error) {
	conn, err := grpc.NewClient(address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create grpc Client: %w", err)
	}

	client := pb.NewStocksServiceClient(conn)

	return &grpcStockService{
		client: client,
		conn:   conn,
	}, nil
}

func (s *grpcStockService) GetStockItemBySKU(ctx context.Context, skuID domain.SkuID) (domain.StockItemBySKU, error) {
	req := &pb.GetStockItemRequest{
		SkuId: uint32(skuID),
	}

	ctx, cancel := context.WithTimeout(ctx, grpcCallTimeOut)
	defer cancel()

	// make the grpc call.
	resp, err := s.client.GetStockItemBySKU(ctx, req)
	if err != nil {
		return domain.StockItemBySKU{}, fmt.Errorf("failed to get stock item via GRPC: %w", err)
	}

	return domain.StockItemBySKU{
		SKuID: domain.SkuID(req.SkuId),
		Name:  resp.Name,
		Price: resp.Price,
		Count: uint16(resp.Count),
	}, nil
}
