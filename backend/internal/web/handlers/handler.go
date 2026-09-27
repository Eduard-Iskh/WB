package handlers

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"wildberies/L0/backend/internal/domain"
	"wildberies/L0/backend/internal/web/handlers/common"

	"github.com/go-chi/chi/v5"
)

type OrderGetter interface {
	GetById(ctx context.Context, id string) (*domain.Order, error)
}

// каждый handler в отдельной папке
func GetOrder(orderService OrderGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		prompt := "get order"

		id := chi.URLParam(r, "id")
		if id == "" {
			common.ErrorResponse(w, fmt.Errorf("%s: empty id", prompt).Error(), http.StatusBadRequest)
			return
		}

		user, err := orderService.GetById(r.Context(), id)
		if err != nil {
			if errors.Is(err, domain.ErrOrderNotFound) {
				common.ErrorResponse(w, fmt.Errorf("%s: %w", prompt, err).Error(), http.StatusNotFound)
				return
			}
			common.ErrorResponse(w, fmt.Errorf("%s: %w", prompt, err).Error(), http.StatusInternalServerError)
			return
		}

		common.SuccessResponse(w, http.StatusOK, map[string]interface{}{"order": user})
	}
}
