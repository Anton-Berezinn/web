package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"projectgrom/internal/cache"
	"projectgrom/internal/token/jwt"
	"time"
)

// Main - hanlder для main page
func (h *Handler) Main(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 1*time.Second)
	defer cancel()
	token := r.Header.Get("Authorization")
	if token == "" {
		w.WriteHeader(401)
		http.Redirect(w, r, "/login", 302)
		return
	}
	err := h.redis.GetValue(ctx, token)
	if err != nil {
		if errors.Is(err, cache.NotFound) {
			w := jwt.ClearToken(w)
			w.WriteHeader(401)
			return
		}
		w.WriteHeader(500)
		return
	}
	data, err := h.productDb.GetAll(ctx)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(203)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, fmt.Sprintf("Failed to encode products: %v", err), http.StatusInternalServerError)
		return
	}
	return
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
}
