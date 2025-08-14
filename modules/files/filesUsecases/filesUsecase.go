package filesUsecases

import "github.com/tonrock01/another-world-shop/config"

type IFilesUsecase interface {
}

type filesUsecase struct {
	cfg config.IConfig
}

func FilesUsecase(cfg config.IConfig) IFilesUsecase {
	return &filesUsecase{
		cfg: cfg,
	}
}
