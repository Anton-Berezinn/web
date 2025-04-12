package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"projectgrom/internal/cache"
	"projectgrom/internal/config"
	"projectgrom/internal/dto"
	"projectgrom/internal/repository/register"
	"projectgrom/internal/token/jwt"
)

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	fmt.Println("here")
	token := r.Header.Get("Authorization")
	if token != "" {
		err := h.redis.Get(token).Err()
		if err != nil {
			if errors.Is(err, cache.NotFound) {
				w.WriteHeader(http.StatusUnauthorized)
				w = jwt.ClearToken(w)
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("tokek was deleted"))
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Println("get data", err)
			return
		}
		http.Redirect(w, r, "/api/main", 301)
		return
	}
	data, err := dto.DataLogin(r.Body)
	if err != nil {
		if errors.Is(err, dto.EmptyDataError) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println("dto login error:", err)
		return
	}
	err = h.register.GetUserLogin(data.Login, data.Password)
	if err != nil {
		if errors.Is(err, register.NoLoginError) {
			w.WriteHeader(400)
			w.Write([]byte("data is wrong"))
			return
		}
		w.WriteHeader(500)
		fmt.Println("error to sql GetUser")
		return
	}
	token, err = jwt.CreateToken(data.Login)
	if err != nil {
		if errors.Is(err, config.EmptyKeyError) {
			w.WriteHeader(500)
			fmt.Println("config is empty")
			return
		}
		w.WriteHeader(500)
		fmt.Println("error to sql CreateToken")
		return
	}
	w.Header().Add("Authorization", token)
	w.WriteHeader(200)
	w.Write([]byte("SUCCESS"))
	return
}
