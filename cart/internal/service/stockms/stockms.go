package stockms

import (
	"bytes"
	"cart/internal/domain"
	"cart/internal/usecase/carts"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const requestTimeout = 10 * time.Second

type stockService struct {
	baseURL    string
	httpClient *http.Client
}

var _ carts.StockService = (*stockService)(nil)

type stockItemResponse struct {
	SkuID uint32 `json:"sku"`
	Name  string `json:"name"`
	Price uint32 `json:"price"`
	Count uint16 `json:"count"`
}

func NewHTTPStockService(baseURL string) *stockService {
	return &stockService{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: requestTimeout,
		},
	}
}

func (s *stockService) GetStockItemBySKU(ctx context.Context, skuID domain.SkuID) (domain.StockItemBySKU, error) {
	reqBody := map[string]uint32{"skuID": uint32(skuID)}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return domain.StockItemBySKU{}, fmt.Errorf("failed to marshal request: %w", err)
	}
	// http request.
	req, err := http.NewRequestWithContext(ctx, "POST", s.baseURL+"/stocks/item/get", bytes.NewBuffer(jsonBody))
	if err != nil {
		return domain.StockItemBySKU{}, fmt.Errorf("failed to create http request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return domain.StockItemBySKU{}, fmt.Errorf("failed to send request to stock service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return domain.StockItemBySKU{}, fmt.Errorf("stock service returned non-200 status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return domain.StockItemBySKU{}, fmt.Errorf("failed to read response body: %w", err)
	}

	var stockItem stockItemResponse
	if err := json.Unmarshal(body, &stockItem); err != nil {
		return domain.StockItemBySKU{}, fmt.Errorf("failed to unmarshal stockItem: %w", err)
	}

	return domain.StockItemBySKU{
		SKuID: domain.SkuID(stockItem.SkuID),
		Name:  stockItem.Name,
		Price: stockItem.Price,
		Count: stockItem.Count,
	}, nil
}
