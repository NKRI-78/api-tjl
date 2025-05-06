package controllers

import (
	"net/http"
	"superapps/entities"
	helper "superapps/helpers"
	"superapps/services"

	"github.com/dgrijalva/jwt-go"
)

func ApplyJobBadges(w http.ResponseWriter, r *http.Request) {
	var dataApplyJobBadges entities.ApplyJobBadges

	tokenHeader := r.Header.Get("Authorization")
	token := helper.DecodeJwt(tokenHeader)

	// Validasi token
	if token == nil || !token.Valid {
		helper.Response(w, http.StatusUnauthorized, true, "Unauthorized or invalid token", nil)
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		helper.Response(w, http.StatusUnauthorized, true, "Invalid token claims", nil)
		return
	}

	userId, ok := claims["id"].(string)
	if !ok {
		helper.Response(w, http.StatusBadRequest, true, "Invalid user ID in token", nil)
		return
	}

	result, err := services.ApplyJobBadges(userId)
	if err != nil {
		helper.Response(w, http.StatusBadRequest, true, err.Error(), nil)
		return
	}

	// Pastikan hasilnya adalah int
	total, ok := result["data"].(int)
	if !ok {
		helper.Response(w, http.StatusInternalServerError, true, "Invalid result data type", nil)
		return
	}

	dataApplyJobBadges.Total = total

	helper.Logger("info", "Apply Job Badges success")
	helper.Response(w, http.StatusOK, false, "Successfully", dataApplyJobBadges)
}
