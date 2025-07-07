package controllers

import (
	"net/http"
	helper "superapps/helpers"
	"superapps/services"

	"github.com/dgrijalva/jwt-go"
)

func AdminApplyJobBadges(w http.ResponseWriter, r *http.Request) {
	tokenHeader := r.Header.Get("Authorization")
	token := helper.DecodeJwt(tokenHeader)
	claims, _ := token.Claims.(jwt.MapClaims)
	BranchId, _ := claims["branch_id"].(string)

	result, err := services.AdminApplyJobBadges(BranchId)
	if err != nil {
		helper.Response(w, http.StatusBadRequest, true, err.Error(), nil)
		return
	}

	badgeData, ok := result["data"].([]map[string]any)
	if !ok {
		helper.Response(w, http.StatusInternalServerError, true, "Invalid result data type", nil)
		return
	}

	helper.Logger("info", "Admin Apply Job Badges success")
	helper.Response(w, http.StatusOK, false, "Successfully", badgeData)
}
