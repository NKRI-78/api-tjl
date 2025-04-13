package controllers

import (
	"encoding/json"
	"net/http"
	helper "superapps/helpers"
	"superapps/models"
	"superapps/services"
)

func CountryUpdate(w http.ResponseWriter, r *http.Request) {

	data := &models.CountryUpdate{}

	err := json.NewDecoder(r.Body).Decode(data)

	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		helper.Response(w, 400, true, "Internal server error ("+err.Error()+")", map[string]any{})
		return
	}

	Name := data.Name
	Currency := data.Currency
	Kurs := data.Kurs
	Info := data.Info

	if Name == "" {
		helper.Logger("error", "In Server: name is required")
		helper.Response(w, 400, true, "name is required", map[string]any{})
		return
	}

	if Currency == "" {
		helper.Logger("error", "In Server: currency is required")
		helper.Response(w, 400, true, "currency is required", map[string]any{})
		return
	}

	if Kurs == "" {
		helper.Logger("error", "In Server: kurs is required")
		helper.Response(w, 400, true, "kurs is required", map[string]any{})
		return
	}

	if Info == "" {
		helper.Logger("error", "In Server: info is required")
		helper.Response(w, 400, true, "info is required", map[string]any{})
		return
	}

	result, err := services.CountryUpdate(data)

	if err != nil {
		helper.Response(w, 400, true, err.Error(), map[string]any{})
		return
	}

	helper.Logger("info", "Update Country success")
	helper.Response(w, http.StatusOK, false, "Successfully", result)
}
