package controllers

import (
	"net/http"
	helper "superapps/helpers"
	"superapps/services"

	"github.com/dgrijalva/jwt-go"
)

func ForumList(w http.ResponseWriter, r *http.Request) {
	page := r.URL.Query().Get("page")
	limit := r.URL.Query().Get("limit")
	search := r.URL.Query().Get("search")

	tokenHeader := r.Header.Get("Authorization")

	token := helper.DecodeJwt(tokenHeader)

	claims, _ := token.Claims.(jwt.MapClaims)

	userId, _ := claims["id"].(string)

	result, err := services.ForumList(userId, search, page, limit)

	if err != nil {
		helper.Response(w, 400, true, err.Error(), map[string]any{})
		return
	}

	helper.Logger("info", "Forum List success")
	helper.ResponseWithPagination(w, http.StatusOK, false, "Successfully",
		result["total"],
		result["per_page"],
		result["prev_page"],
		result["next_page"],
		result["current_page"],
		result["next_url"],
		result["prev_url"],
		result["data"],
	)
}
