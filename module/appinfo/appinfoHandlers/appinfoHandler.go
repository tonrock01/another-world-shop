package appinfoHandlers

import (
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/tonrock01/another-world-shop/config"
	"github.com/tonrock01/another-world-shop/module/appinfo"
	"github.com/tonrock01/another-world-shop/module/appinfo/appinfoUsecases"
	"github.com/tonrock01/another-world-shop/module/entities"
	"github.com/tonrock01/another-world-shop/pkg/anotherworldauth"
)

type appinfoHandlersErrCode string

const (
	generateApiKeyErr appinfoHandlersErrCode = "appinfo-001"
	findCategoryErr   appinfoHandlersErrCode = "appinfo-002"
	addCategoryErr    appinfoHandlersErrCode = "appinfo-003"
	removeCategoryErr appinfoHandlersErrCode = "appinfo-004"
)

type IAppinfoHandler interface {
	GenerateApiKey(c *fiber.Ctx) error
	FindCategory(c *fiber.Ctx) error
	AddCategory(c *fiber.Ctx) error
	RemoveCategory(c *fiber.Ctx) error
}

type appinfoHandler struct {
	cfg             config.IConfig
	appinfoUsercase appinfoUsecases.IAppinfoUsecase
}

func AppinfoHandler(cfg config.IConfig, appinfoUsercase appinfoUsecases.IAppinfoUsecase) IAppinfoHandler {
	return &appinfoHandler{
		cfg:             cfg,
		appinfoUsercase: appinfoUsercase,
	}
}

func (h *appinfoHandler) GenerateApiKey(c *fiber.Ctx) error {
	apiKey, err := anotherworldauth.NewAnotherWorldAuth(
		anotherworldauth.ApiKey,
		h.cfg.Jwt(),
		nil,
	)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(generateApiKeyErr),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(
		fiber.StatusOK,
		&struct {
			Key string `json:"key"`
		}{
			Key: apiKey.SignToken(),
		},
	).Res()
}

func (h *appinfoHandler) FindCategory(c *fiber.Ctx) error {
	req := new(appinfo.CategoryFilter)
	if err := c.QueryParser(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(findCategoryErr),
			err.Error(),
		).Res()
	}

	category, err := h.appinfoUsercase.FindCategory(req)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(findCategoryErr),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(fiber.StatusOK, category).Res()
}

func (h *appinfoHandler) AddCategory(c *fiber.Ctx) error {
	req := make([]*appinfo.Category, 0)
	if err := c.BodyParser(&req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(addCategoryErr),
			err.Error(),
		).Res()
	}
	if len(req) == 0 {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(addCategoryErr),
			"categories request are empty",
		).Res()
	}

	if err := h.appinfoUsercase.InsertCategory(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(addCategoryErr),
			err.Error(),
		).Res()
	}
	return entities.NewResponse(c).Success(fiber.StatusCreated, req).Res()
}

func (h *appinfoHandler) RemoveCategory(c *fiber.Ctx) error {
	categoryId := strings.TrimSpace(c.Params("category_id"))
	categoryIdInt, err := strconv.Atoi(categoryId)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(removeCategoryErr),
			"id type is invalid",
		).Res()
	}
	if categoryIdInt <= 0 {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(removeCategoryErr),
			"id must more than 0",
		).Res()
	}

	if err := h.appinfoUsercase.DeleteCategory(categoryIdInt); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(removeCategoryErr),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(
		fiber.StatusOK,
		&struct {
			CategoryId int `json:"category_id"`
		}{
			CategoryId: categoryIdInt,
		},
	).Res()
}
