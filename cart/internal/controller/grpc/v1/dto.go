package v1

import (
	"cart/internal/domain"
	"cart/pkg/api/cart"
)

func fromGrpcCreateCartItemReqToDomain(req *cart.CreateCartItemRequest) domain.CartItem {
	return domain.CartItem{
		UserID: domain.UserID(req.UserId),
		SkuID:  domain.SkuID(req.SkuId),
		Count:  uint16(req.Count),
	}
}

func fromGrpcDeleteCartItemReqToDomain(req *cart.RemoveCartItemRequest) domain.CartItem {
	return domain.CartItem{
		UserID: domain.UserID(req.UserId),
		SkuID:  domain.SkuID(req.SkuId),
	}
}

func fromListStockItemsDomainToGrpc(cartItemsDomain domain.ListCartItems) *cart.ListCartItemsResponse {
	cartItemsRes := make([]*cart.CartItemResponse, 0, len(cartItemsDomain.Items))

	for _, cartItem := range cartItemsDomain.Items {
		cartItemsRes = append(cartItemsRes, &cart.CartItemResponse{
			SkuId: uint32(cartItem.SKuID),
			Name:  cartItem.Name,
			Count: uint32(cartItem.Count),
			Price: cartItem.Price,
		})
	}

	return &cart.ListCartItemsResponse{
		Items:      cartItemsRes,
		TotalPrice: cartItemsDomain.TotalPrice,
	}
}
