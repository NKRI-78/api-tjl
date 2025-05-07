package controllers

import (
	"net/http"
	helper "superapps/helpers"
	"superapps/services"
)

func ClientList(w http.ResponseWriter, r *http.Request) {

	result, err := services.ClientList()

	if err != nil {
		helper.Response(w, 400, true, err.Error(), map[string]any{})
		return
	}

	helper.Logger("info", "Client List success")
	helper.Response(w, http.StatusOK, false, "Successfully",
		result["data"],
	)
}
