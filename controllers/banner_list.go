package controllers

import (
	"net/http"
	helper "superapps/helpers"
	"superapps/services"
)

func BannerList(w http.ResponseWriter, r *http.Request) {

	result, err := services.BannerList()

	if err != nil {
		helper.Response(w, 400, true, err.Error(), map[string]any{})
		return
	}

	helper.Logger("info", "Get Bannner success")
	helper.Response(w, http.StatusOK, false, "Successfully", result["data"])
}
