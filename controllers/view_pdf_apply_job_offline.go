package controllers

import (
	"net/http"
	helper "superapps/helpers"
	"superapps/services"
)

func ViewPdfApplyJobOffline(w http.ResponseWriter, r *http.Request) {
	applyJobId := r.URL.Query().Get("apply_job_id")

	if applyJobId == "" {
		helper.Logger("error", "In Server: apply_job_id query is required")
		helper.Response(w, 400, true, "apply_job_id query is required", map[string]any{})
		return
	}

	result, err := services.ViewPdfApplyJobOffline(applyJobId)

	if err != nil {
		helper.Response(w, 400, true, err.Error(), map[string]any{})
		return
	}

	helper.Logger("info", "View PDF Departure success")
	helper.Response(w, http.StatusOK, false, "Successfully", result["data"])
}
