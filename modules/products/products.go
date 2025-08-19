package products

import (
	"github.com/tonrock01/another-world-shop/modules/appinfo"
	"github.com/tonrock01/another-world-shop/modules/entities"
)

type Product struct {
	Id          string            `json:"id"`
	Title       string            `json:"title"`
	Description string            `json:"description"`
	Category    *appinfo.Category `json:"category"`
	CreatedAt   string            `json:"created_at"`
	UpdatedAt   string            `json:"updated_at"`
	Price       float64           `json:"price"`
	Images      []*entities.Image `json:"images"`
}

type ProductFilter struct {
	Id     string `query:"id"`
	Search string `query:"search"` // title, description
	*entities.PaginationReq
	*entities.SortReq
}
