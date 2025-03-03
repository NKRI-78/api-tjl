package controllers

import (
	"encoding/json"
	"net/http"
	helper "superapps/helpers"
	"superapps/models"
	"superapps/services"

	"github.com/dgrijalva/jwt-go"
)

func DocumentStore(w http.ResponseWriter, r *http.Request) {

	data := &models.DocumentStore{}

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

	Path := data.Path
	Type := data.Type

	if Path == "" {
		helper.Logger("error", "In Server: path is required")
		helper.Response(w, 400, true, "path is required", map[string]interface{}{})
		return
	}

	if Type == -1 || Type == 0 {
		helper.Logger("error", "In Server: type is required")
		helper.Response(w, 400, true, "type is required", map[string]interface{}{})
	}

	data.Path = Path
	data.Type = Type
	data.UserId = userId

	result, err := services.DocumentStore(data)

	if err != nil {
		helper.Response(w, 400, true, err.Error(), map[string]interface{}{})
		return
	}

	helper.Logger("info", "Store Bannner success")
	helper.Response(w, http.StatusOK, false, "Successfully", result)
}
