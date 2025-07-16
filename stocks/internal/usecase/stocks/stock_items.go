package stocks

import (
	"context"
	"errors"
	"fmt"
	"math"
	"stocks/internal/domain"
	"stocks/internal/kafka"
	"stocks/internal/usecase"
)

//go:generate mkdir -p mock
//go:generate minimock -o ./mock/ -s .go -g
type (
	// SKURepository provides repository methods of SKU service.
	SKURepository interface {
		GetSKUByID(ctx context.Context, skuID domain.SKUID) (domain.SKU, error)
	}

	// StockServiceRepository provides repository methods of stock service.
	StockServiceRepository interface {
		SaveStockItem(ctx context.Context, stockItem domain.StockItem) error
		GetStockItem(ctx context.Context, userID domain.UserID, skuID domain.SKUID) (domain.StockItem, error)
		UpdateStockItem(ctx context.Context, stockItem domain.StockItem) error
		DeleteStockItemFromStorage(ctx context.Context, userID domain.UserID, skuID domain.SKUID) error
		GetStockItemBySku(ct context.Context, skuID domain.SKUID) (domain.StockItem, error)
		ListStockItemsByLocation(ctx context.Context, filter domain.Filter) ([]domain.StockItem, error)
		CountStockItems(ctx context.Context, userID domain.UserID, location string) (uint16, error)
	}
)

type stockServiceUseCase struct {
	SKURepository
	StockServiceRepository
	KafkaProducer kafka.StocksEventProducer
}

var _ usecase.StockServiceUseCase = (*stockServiceUseCase)(nil)

func NewStockServiceUseCase(
	skuRepo SKURepository,
	stockRepo StockServiceRepository,
	kafkaProducer kafka.StocksEventProducer,
) *stockServiceUseCase {
	return &stockServiceUseCase{
		SKURepository:          skuRepo,
		StockServiceRepository: stockRepo,
		KafkaProducer:          kafkaProducer,
	}
}

func (s *stockServiceUseCase) AddStockItem(ctx context.Context, stockItem domain.StockItem) error {
	// check sku exist or not.
	sku, err := s.GetSKUByID(ctx, stockItem.Sku.ID)
	if err != nil {
		return err
	}

	stockItem.Sku = sku

	existingStockItem, err := s.GetStockItem(ctx, stockItem.UserID, stockItem.Sku.ID)
	if err != nil {
		if errors.Is(err, domain.ErrStockItemNotFound) {
			err = s.SaveStockItem(ctx, stockItem)
			if err != nil {
				return err
			}

			s.KafkaProducer.ProduceSKUCreated(ctx, kafka.SKUCreatedAndStockChangedPayload{
				SKU:   fmt.Sprintf("%d", stockItem.Sku.ID),
				Count: stockItem.Count,
				Price: stockItem.Price,
			})
		}

		return err
	}

	stockItem.Count += existingStockItem.Count

	err = s.UpdateStockItem(ctx, stockItem)
	if err != nil {
		return err
	}

	s.KafkaProducer.ProduceStockChanged(ctx, kafka.SKUCreatedAndStockChangedPayload{
		SKU:   fmt.Sprintf("%d", stockItem.Sku.ID),
		Count: stockItem.Count,
		Price: stockItem.Price,
	})

	return nil
}

func (s *stockServiceUseCase) DeleteStockItem(ctx context.Context, userID domain.UserID, skuID domain.SKUID) error {
	return s.DeleteStockItemFromStorage(ctx, userID, skuID)
}

func (s *stockServiceUseCase) GetStockItemBySKU(ctx context.Context, skuID domain.SKUID) (domain.StockItem, error) {
	stockItem, err := s.GetStockItemBySku(ctx, skuID)
	if err != nil {
		return domain.StockItem{}, err
	}

	return stockItem, nil
}

func (s *stockServiceUseCase) ListStockItems(ctx context.Context, filter domain.Filter) (domain.PaginatedResponse[domain.StockItem], error) {
	var paginatedResponse domain.PaginatedResponse[domain.StockItem]
	// count stock items.
	countStockItems, err := s.CountStockItems(ctx, filter.UserID, filter.Location)
	if err != nil {
		return domain.PaginatedResponse[domain.StockItem]{}, err
	}

	paginatedResponse.TotalCount = countStockItems

	listOfStockItems, err := s.ListStockItemsByLocation(ctx, filter)
	if err != nil {
		return domain.PaginatedResponse[domain.StockItem]{}, err
	}

	paginatedResponse.Items = listOfStockItems
	paginatedResponse.PageNumber = int64(math.Ceil(float64(countStockItems) / float64(filter.PageSize)))

	return paginatedResponse, nil
}
