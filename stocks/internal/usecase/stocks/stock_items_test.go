package stocks

import (
	"context"
	"errors"
	"stocks/internal/domain"
	"stocks/internal/usecase/stocks/mock"
	"testing"

	"github.com/gojuno/minimock/v3"
)

// TODO: I MUST FIX TESTS, ONLY STARTED INITIAL STEPS...
func TestStockServiceUseCase_AddStockItem(t *testing.T) {
	t.Parallel()
	_ = minimock.NewController(t)
	ctx := context.Background()

	_ = []struct {
		name          string
		stockItem     domain.StockItem
		skuRepoMock   func(*mock.SKURepositoryMock)
		stockRepoMock func(*mock.StockServiceRepositoryMock)
		wantCount     uint16
		wantErr       bool
		expectedErr   error
	}{
		{
			name: "new item should be saved",
			stockItem: domain.StockItem{
				UserID: 1,
				Sku: domain.SKU{
					ID:   1001,
					Name: "t-shirt",
					Type: "apparel",
				},
				Count:    10,
				Price:    12,
				Location: "Ashgabat",
			},
			skuRepoMock: func(sm *mock.SKURepositoryMock) {
				sm.GetSKUByIDMock.
					Expect(ctx, domain.SKUID(1001)).
					Return(domain.SKU{ID: 1001, Name: "t-shirt", Type: "apparel"}, nil)
			},
			stockRepoMock: func(ssrm *mock.StockServiceRepositoryMock) {
				ssrm.GetStockItemMock.
					Expect(ctx, domain.UserID(1), domain.SKUID(1001)).
					Return(domain.StockItem{}, domain.ErrStockItemNotFound)
				ssrm.SaveStockItemMock.
					Expect(ctx, domain.StockItem{
						UserID:   1,
						Sku:      domain.SKU{ID: 1001, Name: "t-shirt", Type: "apparel"},
						Count:    10,
						Price:    12,
						Location: "Ashgabat",
					}).
					Return(nil)
			},
			wantCount: 10,
			wantErr:   false,
		},
		{
			name: "existing stock item should be updated",
			stockItem: domain.StockItem{
				UserID: 1,
				Sku: domain.SKU{
					ID:   2020,
					Name: "cup",
					Type: "accessory",
				},
				Count:    5,
				Price:    20,
				Location: "Ashgabat",
			},
			skuRepoMock: func(sm *mock.SKURepositoryMock) {
				sm.GetSKUByIDMock.
					Expect(ctx, domain.SKUID(2020)).
					Return(domain.SKU{ID: 2020, Name: "cup", Type: "accessory"}, nil)
			},
			stockRepoMock: func(ssrm *mock.StockServiceRepositoryMock) {
				ssrm.GetStockItemMock.
					Expect(ctx, domain.UserID(1), domain.SKUID(2020)).
					Return(domain.StockItem{
						UserID:   1,
						Sku:      domain.SKU{ID: 2020, Name: "cup", Type: "accessory"},
						Count:    10,
						Price:    12,
						Location: "Ashgabat",
					}, nil)
				ssrm.UpdateStockItemMock.
					Expect(ctx, domain.StockItem{
						UserID:   1,
						Sku:      domain.SKU{ID: 2020, Name: "cup", Type: "accessory"},
						Count:    15,
						Price:    20,
						Location: "Ashgabat",
					}).Return(nil)
			},
			wantCount: 15,
			wantErr:   false,
		},
		{
			name: "sku not found should return error",
			stockItem: domain.StockItem{
				UserID: 1,
				Sku: domain.SKU{
					ID:   1002,
					Name: "t-shirt",
					Type: "apparel",
				},
				Count:    8,
				Price:    30,
				Location: "Ashgabat",
			},
			skuRepoMock: func(sm *mock.SKURepositoryMock) {
				sm.GetSKUByIDMock.
					Expect(ctx, domain.SKUID(1002)).
					Return(domain.SKU{}, domain.ErrSKUNotFound)
			},
			stockRepoMock: func(ssrm *mock.StockServiceRepositoryMock) {},
			wantErr:       true,
			expectedErr:   domain.ErrSKUNotFound,
		},
		{
			name: "get stock item error except not found",
			stockItem: domain.StockItem{
				UserID: 1,
				Sku: domain.SKU{
					ID:   1003,
					Name: "t-shirt",
					Type: "apparel",
				},
				Count:    7,
				Price:    25,
				Location: "Ashgabat",
			},
			skuRepoMock: func(sm *mock.SKURepositoryMock) {
				sm.GetSKUByIDMock.
					Expect(ctx, domain.SKUID(1003)).
					Return(domain.SKU{ID: 1003, Name: "t-shirt", Type: "apparel"}, nil)
			},
			stockRepoMock: func(ssrm *mock.StockServiceRepositoryMock) {
				ssrm.GetStockItemMock.
					Expect(ctx, domain.UserID(1), domain.SKUID(1003)).
					Return(domain.StockItem{}, errors.New("database error"))
			},
			wantErr:     true,
			expectedErr: errors.New("database error"),
		},
		{
			name: "save stock item error",
			stockItem: domain.StockItem{
				UserID: 1,
				Sku: domain.SKU{
					ID:   1002,
					Name: "jacket",
					Type: "apparel",
				},
				Count:    7,
				Price:    25,
				Location: "Ashgabat",
			},
			skuRepoMock: func(sm *mock.SKURepositoryMock) {
				sm.GetSKUByIDMock.
					Expect(ctx, domain.SKUID(1002)).
					Return(domain.SKU{ID: 1002, Name: "t-shirt", Type: "apparel"}, nil)
			},
			stockRepoMock: func(ssrm *mock.StockServiceRepositoryMock) {
				ssrm.GetStockItemMock.
					Expect(ctx, domain.UserID(1), domain.SKUID(1002)).
					Return(domain.StockItem{}, domain.ErrStockItemNotFound)
				ssrm.SaveStockItemMock.
					Expect(ctx, domain.StockItem{
						UserID:   1,
						Sku:      domain.SKU{ID: 1002, Name: "t-shirt", Type: "apparel"},
						Count:    7,
						Price:    25,
						Location: "Ashgabat",
					}).
					Return(errors.New("save stock item failed"))
			},
			wantErr:     true,
			expectedErr: errors.New("save stock item failed"),
		},
	}

	// for _, tt := range tests {
	// 	t.Run(tt.name, func(t *testing.T) {
	// 		skuRepo := mock.NewSKURepositoryMock(ctrl)
	// 		stockRepo := mock.NewStockServiceRepositoryMock(ctrl)
	// 		tt.skuRepoMock(skuRepo)
	// 		tt.stockRepoMock(stockRepo)

	// 		newStockUseCase := NewStockServiceUseCase(skuRepo, stockRepo,)
	// 		err := newStockUseCase.AddStockItem(ctx, tt.stockItem)

	// 		if (err != nil) != tt.wantErr {
	// 			t.Errorf("error=%v, wantErr=%v: AddStockItem()", err, tt.wantErr)
	// 		}
	// 		if tt.wantErr && !errors.Is(err, tt.expectedErr) {
	// 			t.Errorf("error=%v, wantErr=%v: AddStockItem()", err, tt.expectedErr)
	// 		}
	// 	})
	// }
}
