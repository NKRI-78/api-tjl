package controllers

import (
	"net/http"
	helper "superapps/helpers"
	"superapps/services"
)

func JobList(w http.ResponseWriter, r *http.Request) {

	salary := r.URL.Query().Get("salary")
	country := r.URL.Query().Get("country")
	position := r.URL.Query().Get("position")

	result, err := services.JobList(salary, country, position)

	if err != nil {
		helper.Response(w, 400, true, err.Error(), map[string]any{})
		return
	}

	helper.Logger("info", "Fetch Job List success")
	helper.Response(w, http.StatusOK, false, "Successfully", result["data"])
}
