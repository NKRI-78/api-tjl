package controllers

import (
	"encoding/json"
	"net/http"
	"superapps/entities"
	helper "superapps/helpers"
	service "superapps/services"

	"github.com/dgrijalva/jwt-go"
)

func ForumStore(w http.ResponseWriter, r *http.Request) {
	data := &entities.ForumStore{}

	err := json.NewDecoder(r.Body).Decode(data)

	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		helper.Response(w, 400, true, "Internal server error ("+err.Error()+")", map[string]any{})
		return
	}

	tokenHeader := r.Header.Get("Authorization")

	token := helper.DecodeJwt(tokenHeader)

	claims, _ := token.Claims.(jwt.MapClaims)

	userId, _ := claims["id"].(string)

	Title := data.Title
	Desc := data.Caption

	data.UserId = userId

	if Title == "" {
		helper.Logger("error", "In Server: title is required")
		helper.Response(w, 400, true, "title is required", map[string]any{})
		return
	}

	if Desc == "" {
		helper.Logger("error", "In Server: desc is required")
		helper.Response(w, 400, true, "desc is required", map[string]any{})
		return
	}

	result, err := service.ForumStore(data)

	if err != nil {
		helper.Response(w, 400, true, err.Error(), map[string]any{})
		return
	}

	helper.Logger("info", "Forum Store success")
	helper.Response(w, http.StatusOK, false, "Successfully", result["data"])
}
