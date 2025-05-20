package controllers

import (
	"net/http"
	"superapps/entities"
	helper "superapps/helpers"
	"superapps/services"

	"github.com/dgrijalva/jwt-go"
)

func AdminCandidatePassesBadges(w http.ResponseWriter, r *http.Request) {
	var dataAdminApplyJobBadges entities.AdminApplyJobBadges

	tokenHeader := r.Header.Get("Authorization")

	token := helper.DecodeJwt(tokenHeader)

	claims, _ := token.Claims.(jwt.MapClaims)

	BranchId, _ := claims["branch_id"].(string)

	result, err := services.AdminCandidatePassesBadges(BranchId)
	if err != nil {
		helper.Response(w, http.StatusBadRequest, true, err.Error(), nil)
		return
	}

	total, ok := result["data"].(int)
	if !ok {
		helper.Response(w, http.StatusInternalServerError, true, "Invalid result data type", nil)
		return
	}

	dataAdminApplyJobBadges.Total = total

	helper.Logger("info", "Admin Candidate Job Passes Badges success")
	helper.Response(w, http.StatusOK, false, "Successfully", dataAdminApplyJobBadges)
}
