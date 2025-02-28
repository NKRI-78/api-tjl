package controllers

import (
	"encoding/json"
	"net/http"
	helper "superapps/helpers"
	"superapps/models"
	"superapps/services"

	"github.com/dgrijalva/jwt-go"
)

func ProfileUpdate(w http.ResponseWriter, r *http.Request) {

	data := &models.Profile{}

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

	Avatar := data.Avatar
	Fullname := data.Fullname

	data.Id = userId

	if Avatar == "" {
		helper.Logger("error", "In Server: avatar is required")
		helper.Response(w, 400, true, "avatar is required", map[string]interface{}{})
		return
	}

	if Fullname == "" {
		helper.Logger("error", "In Server: fullname is required")
		helper.Response(w, 400, true, "fullname is required", map[string]interface{}{})
		return
	}

	result, err := services.UpdateProfile(data)

	if err != nil {
		helper.Response(w, 400, true, err.Error(), map[string]interface{}{})
		return
	}

	helper.Logger("info", "Store Bannner success")
	helper.Response(w, http.StatusOK, false, "Successfully", result)
}
