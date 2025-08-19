package productsUsecases

import (
	"math"

	"github.com/tonrock01/another-world-shop/modules/entities"
	"github.com/tonrock01/another-world-shop/modules/products"
	"github.com/tonrock01/another-world-shop/modules/products/productsRepositories"
)

type IProductsUsecase interface {
	FindOneProduct(productId string) (*products.Product, error)
	FindProducts(req *products.ProductFilter) *entities.PaginateRes
	AddProduct(req *products.Product) (*products.Product, error)
	UpdateProduct(req *products.Product) (*products.Product, error)
	DeleteProduct(productId string) error
}

type productsUsecase struct {
	productsRepository productsRepositories.IProductsRepository
}

func ProductsUsecase(productsRepository productsRepositories.IProductsRepository) IProductsUsecase {
	return &productsUsecase{
		productsRepository: productsRepository,
	}
}

func (u *productsUsecase) FindOneProduct(productId string) (*products.Product, error) {
	product, err := u.productsRepository.FindOneProduct(productId)
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (u *productsUsecase) FindProducts(req *products.ProductFilter) *entities.PaginateRes {
	product, count := u.productsRepository.FindProducts(req)

	return &entities.PaginateRes{
		Data:      product,
		Page:      req.Page,
		Limit:     req.Limit,
		TotalPage: int(math.Ceil(float64(count) / float64(req.Limit))),
		TotalItem: count,
	}
}

func (u *productsUsecase) AddProduct(req *products.Product) (*products.Product, error) {
	product, err := u.productsRepository.InsertProduct(req)
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (u *productsUsecase) UpdateProduct(req *products.Product) (*products.Product, error) {
	product, err := u.productsRepository.UpdateProduct(req)
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (u *productsUsecase) DeleteProduct(productId string) error {
	if err := u.productsRepository.DeleteProduct(productId); err != nil {
		return err
	}

	return nil
}
