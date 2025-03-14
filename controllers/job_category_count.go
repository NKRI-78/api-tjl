package controllers

import (
	"net/http"
	helper "superapps/helpers"
	"superapps/services"
)

func JobCategoryCount(w http.ResponseWriter, r *http.Request) {
	result, err := services.JobCategoryCount()

	if err != nil {
		helper.Response(w, 400, true, err.Error(), map[string]any{})
		return
	}

	helper.Logger("info", "Fetch Job Category Count success")
	helper.Response(w, http.StatusOK, false, "Successfully", result["data"])
}
