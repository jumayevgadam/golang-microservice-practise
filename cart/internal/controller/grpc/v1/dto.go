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
