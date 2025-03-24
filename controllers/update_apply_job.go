package controllers

import (
	"encoding/json"
	"net/http"
	helper "superapps/helpers"
	"superapps/models"
	"superapps/services"

	"github.com/dgrijalva/jwt-go"
)

func UpdateApplyJob(w http.ResponseWriter, r *http.Request) {

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

	ApplyJobId := data.ApplyJobId
	Link := data.Link
	Schedule := data.Schedule

	UserId := userId

	if data.Status == 1 {
		helper.Logger("error", "In Server: status [IN_PROGRESS] already used")
		helper.Response(w, 400, true, "status [IN_PROGRESS] already used", map[string]any{})
		return
	}

	if ApplyJobId == "" {
		helper.Logger("error", "In Server: apply job id is required")
		helper.Response(w, 400, true, "apply job id is required", map[string]any{})
		return
	}

	if Link == "" {
		helper.Logger("error", "In Server: link is required")
		helper.Response(w, 400, true, "link is required", map[string]any{})
		return
	}

	if Schedule == "" {
		helper.Logger("error", "In Server: schedule is required")
		helper.Response(w, 400, true, "schedule is required", map[string]any{})
		return
	}

	data.UserConfirmId = UserId

	result, err := services.UpdateApplyJob(data)

	if err != nil {
		helper.Response(w, 400, true, err.Error(), map[string]any{})
		return
	}

	helper.Logger("info", "Update Apply Job success")
	helper.Response(w, http.StatusOK, false, "Successfully", result)
}
