package controllers

import (
	"net/http"
	helper "superapps/helpers"
	models "superapps/models"
	"superapps/services"

	"github.com/dgrijalva/jwt-go"
)

func Profile(w http.ResponseWriter, r *http.Request) {
	data := &models.Profile{}

	tokenHeader := r.Header.Get("Authorization")

	token := helper.DecodeJwt(tokenHeader)

	claims, _ := token.Claims.(jwt.MapClaims)

	userId, _ := claims["id"].(string)

	data.Id = userId

	result, err := services.GetProfile(data)

	if err != nil {
		helper.Response(w, 400, true, err.Error(), map[string]any{})
		return
	}

	helper.Response(w, http.StatusOK, false, "Successfully", result["data"])
}
