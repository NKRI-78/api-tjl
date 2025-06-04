package controllers

import (
	"net/http"
	helper "superapps/helpers"
	"superapps/services"
)

func AdminListCandidateImport(w http.ResponseWriter, r *http.Request) {

	result, err := services.AdminListCandidateImport()

	if err != nil {
		helper.Response(w, 400, true, err.Error(), map[string]any{})
		return
	}

	helper.Logger("info", "Get Admin List Candidate Import success")
	helper.Response(w, http.StatusOK, false, "Successfully", result["data"])
}
