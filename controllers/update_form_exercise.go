package controllers

import (
	"encoding/json"
	"net/http"
	helper "superapps/helpers"
	"superapps/models"
	"superapps/services"
)

func UpdateFormExercise(w http.ResponseWriter, r *http.Request) {

	data := &models.FormExercise{}

	err := json.NewDecoder(r.Body).Decode(data)

	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		helper.Response(w, 400, true, "Internal server error ("+err.Error()+")", map[string]interface{}{})
		return
	}

	Id := data.Id
	Institution := data.Institution
	Name := data.Name
	StartYear := data.StartYear
	StartMonth := data.StartMonth
	EndMonth := data.EndMonth
	EndYear := data.EndYear

	if Id == "" {
		helper.Logger("error", "In Server: id is required")
		helper.Response(w, 400, true, "id is required", map[string]interface{}{})
		return
	}

	if Institution == "" {
		helper.Logger("error", "In Server: institution is required")
		helper.Response(w, 400, true, "institution is required", map[string]interface{}{})
		return
	}

	if Name == "" {
		helper.Logger("error", "In Server: name is required")
		helper.Response(w, 400, true, "name is required", map[string]interface{}{})
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

	result, err := services.UpdateFormExercise(data)

	if err != nil {
		helper.Response(w, 400, true, err.Error(), map[string]interface{}{})
		return
	}

	helper.Logger("info", "Update Form Exercise success")
	helper.Response(w, http.StatusOK, false, "Successfully", result)
}
