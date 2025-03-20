package controllers

import (
	"net/http"
	helper "superapps/helpers"
	"superapps/models"
	"superapps/services"

	"github.com/gorilla/mux"
)

func InfoApplyJob(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["id"]

	data := &models.InfoApplyJob{}

	data.Id = id

	result, err := services.InfoApplyJob(data)

	if err != nil {
		helper.Response(w, 400, true, err.Error(), map[string]any{})
		return
	}

	helper.Logger("info", "Get Info Apply Job success")
	helper.Response(w, http.StatusOK, false, "Successfully", result["data"])
}
