package handlers

import (
	"BEWaysBeans/dto"
	"BEWaysBeans/repositories"
	"encoding/json"
	"net/http"

	"github.com/golang-jwt/jwt/v4"
)

type handlerUserTrc struct {
	UserTrcRepository repositories.UserTrcRepository
}

func HandlerUsertrc(UserTrcRepository repositories.UserTrcRepository) *handlerUserTrc {
	return &handlerUserTrc{UserTrcRepository}
}

func (h *handlerUserTrc) FindUserTrc(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userInfo := r.Context().Value("userInfo").(jwt.MapClaims)
	idUser := int(userInfo["id"].(float64))

	// var trips models.Trip
	transaction, err := h.UserTrcRepository.FindUserTrc(idUser)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	// transaction.Attachment = path_file + transaction.Attachment

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Data: transaction}
	json.NewEncoder(w).Encode(response)
}
