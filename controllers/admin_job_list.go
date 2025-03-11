package controllers

import (
	"net/http"
	helper "superapps/helpers"
	"superapps/services"
)

func AdminJobList(w http.ResponseWriter, r *http.Request) {
	result, err := services.AdminJobList()

	if err != nil {
		helper.Response(w, 400, true, err.Error(), map[string]any{})
		return
	}

	helper.Logger("info", "Fetch Admin Job List success")
	helper.Response(w, http.StatusOK, false, "Successfully", result["data"])
}
