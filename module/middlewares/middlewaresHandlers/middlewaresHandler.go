package middlewaresHandlers

import "github.com/tonrock01/another-world-shop/module/middlewares/middlewaresUsecases"

type IMiddlewaresHandler interface {
}

type middlewaresHandler struct {
	middlewaresUsecase middlewaresUsecases.IMiddlewaresUsecase
}

func MiddlewaresRepository(middlewaresUsecase middlewaresUsecases.IMiddlewaresUsecase) IMiddlewaresHandler {
	return &middlewaresHandler{
		middlewaresUsecase: middlewaresUsecase,
	}
}
