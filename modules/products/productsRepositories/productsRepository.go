package productsRepositories

import (
	"github.com/jmoiron/sqlx"
	"github.com/tonrock01/another-world-shop/config"
	"github.com/tonrock01/another-world-shop/modules/files/filesUsecases"
)

type IProductsRepository interface {
}

type productsRepository struct {
	db           *sqlx.DB
	cfg          config.IConfig
	filesUsecase filesUsecases.IFilesUsecase
}

func ProductsRepository(db *sqlx.DB, cfg config.IConfig, filesUsecase filesUsecases.IFilesUsecase) IProductsRepository {
	return &productsRepository{
		db:           db,
		cfg:          cfg,
		filesUsecase: filesUsecase,
	}
}
