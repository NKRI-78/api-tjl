package controllers

import (
	"encoding/json"
	"net/http"
	helper "superapps/helpers"
	"superapps/models"
	"superapps/services"
)

func AssignLanguage(w http.ResponseWriter, r *http.Request) {

	data := &models.FormLanguage{}

	err := json.NewDecoder(r.Body).Decode(data)

	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		helper.Response(w, 400, true, "Internal server error ("+err.Error()+")", map[string]interface{}{})
		return
	}

	Language := data.Language
	Level := data.Level
	UserId := data.UserId

	if Language == "" {
		helper.Logger("error", "In Server: language is required")
		helper.Response(w, 400, true, "language is required", map[string]interface{}{})
		return
	}

	if Level == "" {
		helper.Logger("error", "In Server: level is required")
		helper.Response(w, 400, true, "level is required", map[string]interface{}{})
		return
	}

	if UserId == "" {
		helper.Logger("error", "In Server: user_id is required")
		helper.Response(w, 400, true, "user_id is required", map[string]any{})
		return
	}

	result, err := services.FormLanguage(data)

	if err != nil {
		helper.Response(w, 400, true, err.Error(), map[string]interface{}{})
		return
	}

	helper.Logger("info", "Assign Form Region success")
	helper.Response(w, http.StatusOK, false, "Successfully", result)
}
