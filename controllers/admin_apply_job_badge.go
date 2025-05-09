package controllers

import (
	"net/http"
	"superapps/entities"
	helper "superapps/helpers"
	"superapps/services"
)

func AdminApplyJobBadges(w http.ResponseWriter, r *http.Request) {
	var dataAdminApplyJobBadges entities.AdminApplyJobBadges

	result, err := services.AdminApplyJobBadges()
	if err != nil {
		helper.Response(w, http.StatusBadRequest, true, err.Error(), nil)
		return
	}

	total, ok := result["data"].(int)
	if !ok {
		helper.Response(w, http.StatusInternalServerError, true, "Invalid result data type", nil)
		return
	}

	dataAdminApplyJobBadges.Total = total

	helper.Logger("info", "Admin Apply Job Badges success")
	helper.Response(w, http.StatusOK, false, "Successfully", dataAdminApplyJobBadges)
}
