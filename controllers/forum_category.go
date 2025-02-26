package controllers

import (
	"net/http"
	helper "superapps/helpers"
	"superapps/services"
)

func ForumCategory(w http.ResponseWriter, r *http.Request) {
	result, err := services.ForumCategory()

	if err != nil {
		helper.Response(w, 400, true, err.Error(), map[string]interface{}{})
		return
	}

	helper.Logger("info", "Fetch Forum Category success")
	helper.Response(w, http.StatusOK, false, "Successfully", result["data"])
}
