package v1

import pb "stocks/pkg/api/stocks"

type StockGRPCHandler struct {
	pb.UnimplementedStocksServiceServer
}
