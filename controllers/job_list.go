package controllers

import (
	"net/http"
	helper "superapps/helpers"
	"superapps/services"

	"github.com/dgrijalva/jwt-go"
)

func JobList(w http.ResponseWriter, r *http.Request) {

	salary := r.URL.Query().Get("salary")
	country := r.URL.Query().Get("country")
	position := r.URL.Query().Get("position")

	tokenHeader := r.Header.Get("Authorization")

	token := helper.DecodeJwt(tokenHeader)

	claims, _ := token.Claims.(jwt.MapClaims)

	userId, _ := claims["id"].(string)

	result, err := services.JobList(userId, salary, country, position)

	if err != nil {
		helper.Response(w, 400, true, err.Error(), map[string]any{})
		return
	}

	helper.Logger("info", "Fetch Job List success")
	helper.Response(w, http.StatusOK, false, "Successfully", result["data"])
}
