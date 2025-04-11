package handlers

import (
	"fmt"
	"net/http"
	"projectgrom/web/internal/cache"
	"projectgrom/web/internal/services/products"
)

type Handler struct {
	redis     *cache.RedisCache
	productDb *products.ProductsService
}

func NewHandler(data, dataUser string) (*Handler, error) {
	redis, err := cache.InitRedis()
	if err != nil {
		return nil, err
	}
	product, err := products.InitProductsService(data)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return &Handler{redis: redis, productDb: product}, nil
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {

}
