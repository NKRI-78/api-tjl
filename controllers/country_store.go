package controllers

import (
	"encoding/json"
	"net/http"
	helper "superapps/helpers"
	"superapps/models"
	service "superapps/services"
)

func CountryStore(w http.ResponseWriter, r *http.Request) {

	data := &models.CountryStore{}

	err := json.NewDecoder(r.Body).Decode(data)

	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		helper.Response(w, 400, true, "Internal server error ("+err.Error()+")", map[string]any{})
		return
	}

	Name := data.Name
	Currrency := data.Currency
	Info := data.Info
	Symbol := data.Symbol
	LanguageCode := data.LanguageCode

	if Name == "" {
		helper.Logger("error", "In Server: name is required")
		helper.Response(w, 400, true, "name is required", map[string]any{})
		return
	}

	if Currrency == "" {
		helper.Logger("error", "In Server: currency is required")
		helper.Response(w, 400, true, "currency is required", map[string]any{})
		return
	}

	if Info == "" {
		helper.Logger("error", "In Server: info is required")
		helper.Response(w, 400, true, "info is required", map[string]any{})
		return
	}

	if Symbol == "" {
		helper.Logger("error", "In Server: symbol is required")
		helper.Response(w, 400, true, "symbol is required", map[string]any{})
		return
	}

	if LanguageCode == "" {
		helper.Logger("error", "In Server: language_code is required")
		helper.Response(w, 400, true, "language_code is required", map[string]any{})
		return
	}

	result, err := service.CountryStore(data)

	if err != nil {
		helper.Response(w, 400, true, err.Error(), map[string]any{})
		return
	}

	helper.Logger("info", "Store Job success")
	helper.Response(w, http.StatusOK, false, "Successfully", result)
}
