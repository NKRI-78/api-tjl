package controllers

import (
	"net/http"
	helper "superapps/helpers"
	"superapps/models"
	"superapps/services"

	"github.com/dgrijalva/jwt-go"
)

func ListApplyJob(w http.ResponseWriter, r *http.Request) {

	data := &models.InfoApplyJob{}

	tokenHeader := r.Header.Get("Authorization")

	token := helper.DecodeJwt(tokenHeader)

	claims, _ := token.Claims.(jwt.MapClaims)

	userId, _ := claims["id"].(string)

	data.UserId = userId

	result, err := services.ListInfoApplyJob(data)

	if err != nil {
		helper.Response(w, 400, true, err.Error(), map[string]any{})
		return
	}

	helper.Logger("info", "Get List Info Apply Job success")
	helper.Response(w, http.StatusOK, false, "Successfully", result["data"])
}
