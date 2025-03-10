package controllers

import (
	"net/http"
	helper "superapps/helpers"
	"superapps/models"
	"superapps/services"

	"github.com/gorilla/mux"
)

func District(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["city_id"]

	data := &models.District{}

	data.RegencyId = id

	result, err := services.District(data)

	if err != nil {
		helper.Response(w, 400, true, err.Error(), map[string]any{})
		return
	}

	helper.Logger("info", "Get District success")
	helper.Response(w, http.StatusOK, false, "Successfully", result["data"])
}
