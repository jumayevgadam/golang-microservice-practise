package v1

import (
	"cart/internal/domain"
	"cart/pkg/api/cart"
	helper "cart/pkg/httphelper"
)

func fromGrpcCreateCartItemReqToDomain(req *cart.CreateCartItemRequest) (domain.CartItem, error) {
	createCartItemReq := CreateCartItemRequest{
		UserID: req.UserId,
		SkuID:  req.SkuId,
		Count:  uint16(req.Count),
	}

	if err := helper.ValidateRequest(&createCartItemReq); err != nil {
		return domain.CartItem{}, err
	}

	return createCartItemReq.ToDomain(), nil
}

func fromGrpcDeleteCartItemReqToDomain(req *cart.RemoveCartItemRequest) (domain.CartItem, error) {
	deleteCartItemReq := DeleteCartItemRequest{
		UserID: req.UserId,
		SkuID:  req.SkuId,
	}

	if err := helper.ValidateRequest(&deleteCartItemReq); err != nil {
		return domain.CartItem{}, err
	}

	return domain.CartItem{
		UserID: domain.UserID(deleteCartItemReq.UserID),
		SkuID:  domain.SkuID(deleteCartItemReq.SkuID),
	}, nil
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
