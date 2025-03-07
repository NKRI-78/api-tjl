package controllers

import (
	"net/http"
	helper "superapps/helpers"
	models "superapps/models"
	"superapps/services"

	"github.com/gorilla/mux"
)

func JobDetail(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	data := &models.Job{}

	data.Id = id

	result, err := services.JobDetail(data)

	if err != nil {
		helper.Response(w, 400, true, err.Error(), map[string]any{})
		return
	}

	helper.Logger("info", "Job Detail success")
	helper.Response(w, http.StatusOK, false, "Successfully",
		result["data"],
	)
}
