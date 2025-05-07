package controllers

import (
	"net/http"
	"superapps/entities"
	helper "superapps/helpers"
	"superapps/services"

	"github.com/dgrijalva/jwt-go"
)

func InboxBadge(w http.ResponseWriter, r *http.Request) {
	var dataInboxBadge entities.InboxBadge

	tokenHeader := r.Header.Get("Authorization")

	token := helper.DecodeJwt(tokenHeader)

	claims, _ := token.Claims.(jwt.MapClaims)

	userId, _ := claims["id"].(string)

	result, err := services.InboxBadge(userId)

	if err != nil {
		helper.Response(w, 400, true, err.Error(), map[string]any{})
		return
	}

	dataInboxBadge.Total = result["data"].(int)

	helper.Logger("info", "Get Inbox List success")
	helper.Response(w, http.StatusOK, false, "Successfully", dataInboxBadge)
}
