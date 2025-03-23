package controllers

import (
	"encoding/json"
	"net/http"
	helper "superapps/helpers"
	"superapps/models"
	"superapps/services"
)

func AssignDocumentApplyJob(w http.ResponseWriter, r *http.Request) {

	data := &models.AssignDocumentApplyJob{}

	err := json.NewDecoder(r.Body).Decode(data)

	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		helper.Response(w, 400, true, "Internal server error ("+err.Error()+")", map[string]any{})
		return
	}

	ApplyJobId := data.ApplyJobId
	Path := data.Path

	if ApplyJobId == "" {
		helper.Logger("error", "In Server: apply_job_id is required")
		helper.Response(w, 400, true, "job_id is required", map[string]any{})
		return
	}

	if Path == "" {
		helper.Logger("error", "In Server: path is required")
		helper.Response(w, 400, true, "path is required", map[string]any{})
		return
	}

	result, err := services.AssignDocumentApplyJob(data)

	if err != nil {
		helper.Response(w, 400, true, err.Error(), map[string]any{})
		return
	}

	helper.Logger("info", "Apply Job success")
	helper.Response(w, http.StatusOK, false, "Successfully", result["data"])
}
