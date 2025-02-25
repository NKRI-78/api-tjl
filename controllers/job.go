package controllers

import (
	"net/http"
	helper "superapps/helpers"
	"superapps/services"
)

func JobCategories(w http.ResponseWriter, r *http.Request) {
	result, err := services.GetJobCategory()

	if err != nil {
		helper.Response(w, 400, true, err.Error(), map[string]interface{}{})
		return
	}

	helper.Response(w, http.StatusOK, false, "Successfully", result["data"])
}
