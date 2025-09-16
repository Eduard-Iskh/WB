package handlers

import (
	"fmt"
	"log"
	"net/http"
	"wildberies/L0/backend/handlers/common"
	"wildberies/L0/backend/internal/app"

	"github.com/go-chi/chi"
)

// каждый handler в отдельной папке
func GetOrder(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		prompt := "получение информации о пользователе"

		log.Println("request", r)

		id := chi.URLParam(r, "id")
		if id == "" {
			common.ErrorResponse(w, fmt.Errorf("%s: пустой id", prompt).Error(), http.StatusBadRequest)
			return
		}

		// вызвать сервис GetById
		// user, err := app.UserService.GetById(r.Context(), idUuid)
		// if err != nil {
		// 	common.ErrorResponse(w, fmt.Errorf("%s: %w", prompt, err).Error(), http.StatusInternalServerError)
		// 	return
		// }

		//common.SuccessResponse(w, http.StatusOK, map[string]interface{}{"user": common.ToUserTransport(user)})

		common.SuccessResponse(w, http.StatusOK, map[string]interface{}{"user": id})
	}
}
