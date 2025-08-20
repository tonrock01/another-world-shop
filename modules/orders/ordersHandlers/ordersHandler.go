package ordersHandlers

import (
	"github.com/tonrock01/another-world-shop/config"
	"github.com/tonrock01/another-world-shop/modules/orders/ordersUsecases"
)

type IOrdersHandler interface {
}

type ordersHandler struct {
	cfg            config.IConfig
	ordersUsecases ordersUsecases.IOrdersUsecase
}

func OrdersHandler(cfg config.IConfig, ordersUsecases ordersUsecases.IOrdersUsecase) IOrdersHandler {
	return &ordersHandler{
		cfg:            cfg,
		ordersUsecases: ordersUsecases,
	}
}
