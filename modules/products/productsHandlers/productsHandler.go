package productsHandlers

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/tonrock01/another-world-shop/config"
	"github.com/tonrock01/another-world-shop/modules/appinfo"
	"github.com/tonrock01/another-world-shop/modules/entities"
	"github.com/tonrock01/another-world-shop/modules/files"
	"github.com/tonrock01/another-world-shop/modules/files/filesUsecases"
	"github.com/tonrock01/another-world-shop/modules/products"
	"github.com/tonrock01/another-world-shop/modules/products/productsUsecases"
)

type productsHandlersErrCode string

const (
	findOneProductErr productsHandlersErrCode = "products-001"
	findProductsErr   productsHandlersErrCode = "products-002"
	insertProductsErr productsHandlersErrCode = "products-003"
	updateProductsErr productsHandlersErrCode = "products-004"
	deleteProductsErr productsHandlersErrCode = "products-005"
)

type IProductsHandler interface {
	FindOneProduct(c *fiber.Ctx) error
	FindProducts(c *fiber.Ctx) error
	AddProduct(c *fiber.Ctx) error
	UpdateProduct(c *fiber.Ctx) error
	DeleteProduct(c *fiber.Ctx) error
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

func (h *productsHandler) FindOneProduct(c *fiber.Ctx) error {
	productId := strings.TrimSpace(c.Params("product_id"))

	product, err := h.productsUsecases.FindOneProduct(productId)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(findOneProductErr),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(fiber.StatusOK, product).Res()
}

func (h *productsHandler) FindProducts(c *fiber.Ctx) error {
	req := &products.ProductFilter{
		PaginationReq: &entities.PaginationReq{},
		SortReq:       &entities.SortReq{},
	}

	if err := c.QueryParser(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(findProductsErr),
			err.Error(),
		).Res()
	}

	if req.Page < 1 {
		req.Page = 1
	}
	if req.Limit < 5 {
		req.Limit = 5
	}

	if req.OrderBy == "" {
		req.OrderBy = "title"
	}
	if req.Sort == "" {
		req.Sort = "ASC"
	}

	products := h.productsUsecases.FindProducts(req)
	return entities.NewResponse(c).Success(fiber.StatusOK, products).Res()
}

func (h *productsHandler) AddProduct(c *fiber.Ctx) error {
	req := &products.Product{
		Category: &appinfo.Category{},
		Images:   make([]*entities.Image, 0),
	}
	if err := c.BodyParser(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(insertProductsErr),
			err.Error(),
		).Res()
	}

	if req.Category.Id <= 0 {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(insertProductsErr),
			"category id is invalid",
		).Res()
	}

	product, err := h.productsUsecases.AddProduct(req)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(insertProductsErr),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(fiber.StatusCreated, product).Res()
}

func (h *productsHandler) UpdateProduct(c *fiber.Ctx) error {
	productId := strings.TrimSpace(c.Params("product_id"))

	req := &products.Product{
		Images:   make([]*entities.Image, 0),
		Category: &appinfo.Category{},
	}

	if err := c.BodyParser(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(updateProductsErr),
			err.Error(),
		).Res()
	}

	req.Id = productId

	product, err := h.productsUsecases.UpdateProduct(req)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(updateProductsErr),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(fiber.StatusOK, product).Res()
}

func (h *productsHandler) DeleteProduct(c *fiber.Ctx) error {
	productId := strings.TrimSpace(c.Params("product_id"))

	product, err := h.productsUsecases.FindOneProduct(productId)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(deleteProductsErr),
			err.Error(),
		).Res()
	}

	deleteFileReq := make([]*files.DeleteFileReq, 0)

	for _, p := range product.Images {
		deleteFileReq = append(deleteFileReq, &files.DeleteFileReq{
			Destination: fmt.Sprintf("images/products/%s", p.FileName),
		})
	}

	if err := h.filesUsecases.DeleteFileOnGCP(deleteFileReq); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(deleteProductsErr),
			err.Error(),
		).Res()
	}

	if err := h.productsUsecases.DeleteProduct(product.Id); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(deleteProductsErr),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(fiber.StatusNoContent, nil).Res()
}
