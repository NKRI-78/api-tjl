package controllers

import (
	"net/http"
	helper "superapps/helpers"
	"superapps/models"
	"superapps/services"

	"github.com/gorilla/mux"
)

func City(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["province_id"]

	data := &models.City{}

	data.ProvinceId = id

	result, err := services.City(data)

	if err != nil {
		helper.Response(w, 400, true, err.Error(), map[string]any{})
		return
	}

	helper.Logger("info", "Get City success")
	helper.Response(w, http.StatusOK, false, "Successfully", result["data"])
}
