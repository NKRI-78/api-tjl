package controllers

import (
	"encoding/json"
	"net/http"
	"superapps/entities"
	helper "superapps/helpers"
	"superapps/services"

	"github.com/dgrijalva/jwt-go"
)

func EventStore(w http.ResponseWriter, r *http.Request) {

	data := &entities.EventStore{}

	err := json.NewDecoder(r.Body).Decode(data)

	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		helper.Response(w, 400, true, "Internal server error ("+err.Error()+")", map[string]any{})
		return
	}

	tokenHeader := r.Header.Get("Authorization")

	token := helper.DecodeJwt(tokenHeader)

	claims, _ := token.Claims.(jwt.MapClaims)

	UserId, _ := claims["id"].(string)

	Title := data.Title
	Caption := data.Caption

	data.Title = Title
	data.Caption = Caption
	data.UserId = UserId

	if Title == "" {
		helper.Logger("error", "In Server: title is required")
		helper.Response(w, 400, true, "title is required", map[string]any{})
		return
	}

	if Caption == "" {
		helper.Logger("error", "In Server: caption is required")
		helper.Response(w, 400, true, "caption is required", map[string]any{})
		return
	}

	result, err := services.EventStore(data)

	if err != nil {
		helper.Response(w, 400, true, err.Error(), map[string]any{})
		return
	}

	Id := result["id"].(int64)

	helper.Logger("info", "Store Bannner success")
	helper.Response(w, http.StatusOK, false, "Successfully", Id)
}
