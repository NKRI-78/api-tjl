package controllers

import (
	"encoding/json"
	"net/http"
	helper "superapps/helpers"
	"superapps/models"
	"superapps/services"

	"github.com/dgrijalva/jwt-go"
)

func FormBiodata(w http.ResponseWriter, r *http.Request) {

	data := &models.FormBiodata{}

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

	Place := data.Place
	Birthdate := data.Birthdate
	Gender := data.Gender
	Height := data.Height
	Weight := data.Weight
	Religion := data.Religion
	Status := data.Status
	UserId := userId

	if Place == "" {
		helper.Logger("error", "In Server: place is required")
		helper.Response(w, 400, true, "place is required", map[string]any{})
		return
	}

	if Birthdate == "" {
		helper.Logger("error", "In Server: birthdate is required")
		helper.Response(w, 400, true, "birthdate is required", map[string]any{})
		return
	}

	if Gender == "" {
		helper.Logger("error", "In Server: gender is required")
		helper.Response(w, 400, true, "gender is required", map[string]any{})
		return
	}

	if Height == "" {
		helper.Logger("error", "In Server: height is required")
		helper.Response(w, 400, true, "height is required", map[string]any{})
		return
	}

	if Weight == "" {
		helper.Logger("error", "In Server: weight is required")
		helper.Response(w, 400, true, "weight is required", map[string]any{})
		return
	}

	if Religion == "" {
		helper.Logger("error", "In Server: religion is required")
		helper.Response(w, 400, true, "religion is required", map[string]any{})
		return
	}

	if Status == "" {
		helper.Logger("error", "In Server: status is required")
		helper.Response(w, 400, true, "status is required", map[string]any{})
		return
	}

	data.UserId = UserId

	result, err := services.FormBiodata(data)

	if err != nil {
		helper.Response(w, 400, true, err.Error(), map[string]any{})
		return
	}

	helper.Logger("info", "Form Biodata success")
	helper.Response(w, http.StatusOK, false, "Successfully", result)
}
