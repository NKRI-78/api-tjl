package controllers

import (
	"encoding/json"
	"net/http"
	helper "superapps/helpers"
	"superapps/models"
	"superapps/services"

	"github.com/dgrijalva/jwt-go"
)

func FormLanguage(w http.ResponseWriter, r *http.Request) {

	data := &models.FormLanguage{}

	err := json.NewDecoder(r.Body).Decode(data)

	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		helper.Response(w, 400, true, "Internal server error ("+err.Error()+")", map[string]interface{}{})
		return
	}

	tokenHeader := r.Header.Get("Authorization")

	token := helper.DecodeJwt(tokenHeader)

	claims, _ := token.Claims.(jwt.MapClaims)

	userId, _ := claims["id"].(string)

	Language := data.Language
	Level := data.Level

	UserId := userId

	if Language == "" {
		helper.Logger("error", "In Server: province_id is required")
		helper.Response(w, 400, true, "province_id is required", map[string]interface{}{})
		return
	}

	if Level == "" {
		helper.Logger("error", "In Server: city_id is required")
		helper.Response(w, 400, true, "city_id is required", map[string]interface{}{})
		return
	}

	data.UserId = UserId

	result, err := services.FormLanguage(data)

	if err != nil {
		helper.Response(w, 400, true, err.Error(), map[string]interface{}{})
		return
	}

	helper.Logger("info", "Form Region success")
	helper.Response(w, http.StatusOK, false, "Successfully", result)
}
