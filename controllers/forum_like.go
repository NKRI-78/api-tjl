package controllers

import (
	"encoding/json"
	"net/http"
	"superapps/entities"
	helper "superapps/helpers"
	service "superapps/services"

	"github.com/dgrijalva/jwt-go"
)

func ForumLike(w http.ResponseWriter, r *http.Request) {
	data := &entities.ForumStoreLike{}

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

	ForumId := data.ForumId

	data.UserId = userId

	if ForumId == "" {
		helper.Logger("error", "In Server: forum_id is required")
		helper.Response(w, 400, true, "forum_id is required", map[string]any{})
		return
	}

	result, err := service.ForumStoreLike(data)

	if err != nil {
		helper.Response(w, 400, true, err.Error(), map[string]any{})
		return
	}

	helper.Logger("info", "Forum Store success")
	helper.Response(w, http.StatusOK, false, "Successfully", result["message"])
}
