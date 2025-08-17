package productsHandlers

import (
	"github.com/tonrock01/another-world-shop/config"
	"github.com/tonrock01/another-world-shop/modules/files/filesUsecases"
	"github.com/tonrock01/another-world-shop/modules/products/productsUsecases"
)

type IProductsHandler interface {
}

type productsHandler struct {
	cfg              config.IConfig
	productsUsecases productsUsecases.IProductsUsecase
	filesUsecases    filesUsecases.IFilesUsecase
}

func ProductsHandler(cfg config.IConfig, productsUsecases productsUsecases.IProductsUsecase, filesUsecases filesUsecases.IFilesUsecase) IProductsHandler {
	return &productsHandler{
		cfg:              cfg,
		productsUsecases: productsUsecases,
		filesUsecases:    filesUsecases,
	}
}
