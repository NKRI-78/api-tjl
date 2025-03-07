package controllers

import (
	"encoding/json"
	"net/http"
	helper "superapps/helpers"
	"superapps/models"
	"superapps/services"

	"github.com/dgrijalva/jwt-go"
)

func FormEducation(w http.ResponseWriter, r *http.Request) {

	data := &models.FormEducation{}

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

	EducationLevel := data.EducationLevel
	Major := data.Major
	SchoolOrCollege := data.SchoolOrCollege
	StartYear := data.StartYear
	StartMonth := data.StartMonth
	EndMonth := data.EndMonth
	EndYear := data.EndYear

	UserId := userId

	if EducationLevel == "" {
		helper.Logger("error", "In Server: education_level is required")
		helper.Response(w, 400, true, "education_level is required", map[string]interface{}{})
		return
	}

	if Major == "" {
		helper.Logger("error", "In Server: major is required")
		helper.Response(w, 400, true, "major is required", map[string]interface{}{})
		return
	}

	if SchoolOrCollege == "" {
		helper.Logger("error", "In Server: school_or_college is required")
		helper.Response(w, 400, true, "school_or_college is required", map[string]interface{}{})
		return
	}

	if StartYear == "" {
		helper.Logger("error", "In Server: start_year is required")
		helper.Response(w, 400, true, "start_year is required", map[string]interface{}{})
		return
	}

	if StartMonth == "" {
		helper.Logger("error", "In Server: start_month is required")
		helper.Response(w, 400, true, "start_month is required", map[string]interface{}{})
		return
	}

	if EndMonth == "" {
		helper.Logger("error", "In Server: end_month is required")
		helper.Response(w, 400, true, "end_month is required", map[string]interface{}{})
		return
	}

	if EndYear == "" {
		helper.Logger("error", "In Server: end_year is required")
		helper.Response(w, 400, true, "end_year is required", map[string]interface{}{})
		return
	}

	data.UserId = UserId

	result, err := services.FormEducation(data)

	if err != nil {
		helper.Response(w, 400, true, err.Error(), map[string]interface{}{})
		return
	}

	helper.Logger("info", "Form Education success")
	helper.Response(w, http.StatusOK, false, "Successfully", result)
}
