package controllers

import (
	"net/http"
	helper "superapps/helpers"
	"superapps/services"

	"github.com/dgrijalva/jwt-go"
)

func ViewPdfDeparture(w http.ResponseWriter, r *http.Request) {
	applyJobId := r.URL.Query().Get("apply_job_id")

	if applyJobId == "" {
		helper.Logger("error", "In Server: apply_job_id query is required")
		helper.Response(w, 400, true, "apply_job_id query is required", map[string]any{})
		return
	}

	tokenHeader := r.Header.Get("Authorization")

	token := helper.DecodeJwt(tokenHeader)

	claims, _ := token.Claims.(jwt.MapClaims)

	userId, _ := claims["id"].(string)

	result, err := services.ViewPdfDeparture(userId, applyJobId)

	if err != nil {
		helper.Response(w, 400, true, err.Error(), map[string]any{})
		return
	}

	helper.Logger("info", "View PDF Departure success")
	helper.Response(w, http.StatusOK, false, "Successfully", result["data"])
}
