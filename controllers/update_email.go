package controllers

import (
	"encoding/json"
	"net/http"
	helper "superapps/helpers"
	"superapps/models"
	"superapps/services"
)

func UpdateEmail(w http.ResponseWriter, r *http.Request) {

	data := &models.UpdateEmail{}

	err := json.NewDecoder(r.Body).Decode(data)

	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		helper.Response(w, 400, true, "Internal server error ("+err.Error()+")", map[string]any{})
		return
	}

	OldEmail := data.OldEmail
	NewEmail := data.NewEmail

	if OldEmail == "" {
		helper.Logger("error", "In Server: old_email is required")
		helper.Response(w, 400, true, "old_email is required", map[string]any{})
		return
	}

	if NewEmail == "" {
		helper.Logger("error", "In Server: new_email id is required")
		helper.Response(w, 400, true, "new_email id is required", map[string]any{})
		return
	}

	result, err := services.UpdateEmail(data)

	if err != nil {
		helper.Response(w, 400, true, err.Error(), map[string]any{})
		return
	}

	helper.Logger("info", "Update Apply Job success")
	helper.Response(w, http.StatusOK, false, "Successfully", result)
}
