package controllers

import (
	"net/http"
	helper "superapps/helpers"
	"superapps/models"
	"superapps/services"

	"github.com/gorilla/mux"
)

func Subdistrict(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["district_id"]

	data := &models.Subdistrict{}

	data.DistrictId = id

	result, err := services.Subdistrict(data)

	if err != nil {
		helper.Response(w, 400, true, err.Error(), map[string]any{})
		return
	}

	helper.Logger("info", "Get Subdistrict success")
	helper.Response(w, http.StatusOK, false, "Successfully", result["data"])
}
