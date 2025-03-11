package controllers

import (
	"encoding/json"
	"net/http"
	helper "superapps/helpers"
	"superapps/models"
	"superapps/services"

	"github.com/dgrijalva/jwt-go"
)

func ApplyJob(w http.ResponseWriter, r *http.Request) {

	data := &models.ApplyJob{}

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

	data.UserId = userId
	JobId := data.JobId
	Status := data.Status

	if JobId == "" {
		helper.Logger("error", "In Server: job_id is required")
		helper.Response(w, 400, true, "job_id is required", map[string]any{})
		return
	}

	if Status == "0" || Status == "" {
		helper.Logger("error", "In Server: status is required")
		helper.Response(w, 400, true, "status is required", map[string]any{})
		return
	}

	result, err := services.ApplyJob(data)

	if err != nil {
		helper.Response(w, 400, true, err.Error(), map[string]any{})
		return
	}

	helper.Logger("info", "Apply Job success")
	helper.Response(w, http.StatusOK, false, "Successfully", result["data"])
}
