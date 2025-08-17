package productsUsecases

import (
	"github.com/tonrock01/another-world-shop/modules/products/productsRepositories"
)

type IProductsUsecase interface {
}

type productsUsecase struct {
	productsRepository productsRepositories.IProductsRepository
}

func ProductsUsecase(productsRepository productsRepositories.IProductsRepository) IProductsUsecase {
	return &productsUsecase{
		productsRepository: productsRepository,
	}
}
