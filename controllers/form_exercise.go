package controllers

import (
	"encoding/json"
	"net/http"
	helper "superapps/helpers"
	"superapps/models"
	"superapps/services"

	"github.com/dgrijalva/jwt-go"
)

func FormExercise(w http.ResponseWriter, r *http.Request) {

	data := &models.FormExercise{}

	err := json.NewDecoder(r.Body).Decode(data)

	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		helper.Response(w, 400, true, "Internal server error ("+err.Error()+")", map[string]interface{}{})
		return
	}

	tokenHeader := r.Header.Get("Authorization")

	token := helper.DecodeJwt(tokenHeader)

	claims, _ := token.Claims.(jwt.MapClaims)

	userId, _ := claims["id"].(string)

	Name := data.Name
	Institution := data.Institution
	StartYear := data.StartYear
	StartMonth := data.StartMonth
	EndMonth := data.EndMonth
	EndYear := data.EndYear

	UserId := userId

	if Name == "" {
		helper.Logger("error", "In Server: name is required")
		helper.Response(w, 400, true, "name is required", map[string]any{})
		return
	}

	if Institution == "" {
		helper.Logger("error", "In Server: institution is required")
		helper.Response(w, 400, true, "institution is required", map[string]any{})
		return
	}

	if StartYear == "" {
		helper.Logger("error", "In Server: start_year is required")
		helper.Response(w, 400, true, "start_year is required", map[string]any{})
		return
	}

	if StartMonth == "" {
		helper.Logger("error", "In Server: start_month is required")
		helper.Response(w, 400, true, "start_month is required", map[string]any{})
		return
	}

	if EndMonth == "" {
		helper.Logger("error", "In Server: end_month is required")
		helper.Response(w, 400, true, "end_month is required", map[string]any{})
		return
	}

	if EndYear == "" {
		helper.Logger("error", "In Server: end_year is required")
		helper.Response(w, 400, true, "end_year is required", map[string]any{})
		return
	}

	data.UserId = UserId

	result, err := services.FormExercise(data)

	if err != nil {
		helper.Response(w, 400, true, err.Error(), map[string]any{})
		return
	}

	helper.Logger("info", "Form Exercise success")
	helper.Response(w, http.StatusOK, false, "Successfully", result)
}
