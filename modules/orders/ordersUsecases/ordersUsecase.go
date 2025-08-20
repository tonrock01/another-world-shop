package ordersUsecases

import (
	"github.com/tonrock01/another-world-shop/modules/orders/ordersRepositories"
	"github.com/tonrock01/another-world-shop/modules/products/productsRepositories"
)

type IOrdersUsecase interface {
}

type ordersUsecase struct {
	ordersRepository   ordersRepositories.IOrdersRepository
	productsRepository productsRepositories.IProductsRepository
}

func OrdersUsecase(ordersRepository ordersRepositories.IOrdersRepository, productsRepository productsRepositories.IProductsRepository) IOrdersUsecase {
	return &ordersUsecase{
		ordersRepository:   ordersRepository,
		productsRepository: productsRepository,
	}
}
