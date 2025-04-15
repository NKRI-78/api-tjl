package controllers

import (
	"encoding/json"
	"net/http"
	"superapps/entities"
	helper "superapps/helpers"
	"superapps/services"
)

func CompanyUpdate(w http.ResponseWriter, r *http.Request) {

	data := &entities.CompanyUpdate{}

	err := json.NewDecoder(r.Body).Decode(data)

	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		helper.Response(w, 400, true, "Internal server error ("+err.Error()+")", map[string]any{})
		return
	}

	Logo := data.Logo
	Name := data.Name

	if Logo == "" {
		helper.Logger("error", "In Server: logo is required")
		helper.Response(w, 400, true, "logo is required", map[string]any{})
		return
	}

	if Name == "" {
		helper.Logger("error", "In Server: name is required")
		helper.Response(w, 400, true, "name is required", map[string]any{})
		return
	}

	result, err := services.CompanyUpdate(data)

	if err != nil {
		helper.Response(w, 400, true, err.Error(), map[string]any{})
		return
	}

	helper.Logger("info", "Update Country success")
	helper.Response(w, http.StatusOK, false, "Successfully", result)
}
