package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"projectgrom/internal/cache"
	"projectgrom/internal/config"
	"projectgrom/internal/dto"
	storag "projectgrom/internal/repository/register"
	"projectgrom/internal/services/products"
	"projectgrom/internal/services/register"
	"projectgrom/internal/token/jwt"
)

type Handler struct {
	redis     *cache.RedisCache
	productDb *products.ProductsService
	register  *register.RegisterService
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
	regisServe, err := register.NewRegisterService(dataUser)
	return &Handler{redis: redis, productDb: product, register: regisServe}, nil
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	data, err := dto.DataRequest(r.Body)
	if err != nil {
		if errors.Is(err, dto.EmptyDataError) {
			w.WriteHeader(401)
			return
		}
		fmt.Println("error in requestdata err:", err)
		w.WriteHeader(500)
		return
	}
	err = h.register.Add(data.Firstname, data.Lastname, data.Login, data.Password)
	if err != nil {
		if errors.Is(err, storag.LenPassError) {
			w.WriteHeader(401)
			w.Write([]byte("Len pass is more than 6"))
			return
		}
	}
	token, err := jwt.CreateToken(data.Login)
	if err != nil {
		if errors.Is(err, config.EmptyKeyError) {
			fmt.Println("empty key")
		}
		w.WriteHeader(500)
		return
	}
	err = h.redis.Add(token, data.Login)
	if err != nil {
		fmt.Println("error to redis add")
		w.WriteHeader(500)
		return
	}
	w.Header().Add("Authorization", token)
	w.WriteHeader(201)
	w.Write([]byte("SUCCESS"))
	return
}
