package controllers

import (
	"net/http"
	helper "superapps/helpers"
	"superapps/services"

	"github.com/dgrijalva/jwt-go"
)

func AdminListApplyJob(w http.ResponseWriter, r *http.Request) {

	Filter := r.URL.Query().Get("filter")

	tokenHeader := r.Header.Get("Authorization")

	token := helper.DecodeJwt(tokenHeader)

	claims, _ := token.Claims.(jwt.MapClaims)

	BranchId, _ := claims["branch_id"].(string)

	result, err := services.AdminListApplyJob(BranchId, Filter)

	if err != nil {
		helper.Response(w, 400, true, err.Error(), map[string]any{})
		return
	}

	helper.Logger("info", "Fetch Admin Job List success")
	helper.Response(w, http.StatusOK, false, "Successfully", result["data"])
}
