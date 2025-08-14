package filesHandlers

import (
	"github.com/tonrock01/another-world-shop/config"
	"github.com/tonrock01/another-world-shop/modules/files/filesUsecases"
)

type IFilesHandler interface {
}

type filesHandler struct {
	cfg          config.IConfig
	filesUsecase filesUsecases.IFilesUsecase
}

func FilesHandler(cfg config.IConfig, filesUsecase filesUsecases.IFilesUsecase) IFilesHandler {
	return &filesHandler{
		cfg:          cfg,
		filesUsecase: filesUsecase,
	}
}
