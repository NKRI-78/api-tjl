package controllers

import (
	"encoding/json"
	"net/http"
	helper "superapps/helpers"
	"superapps/models"
	service "superapps/services"
)

func JobFavourite(w http.ResponseWriter, r *http.Request) {

	data := &models.JobFavourite{}

	err := json.NewDecoder(r.Body).Decode(data)

	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		helper.Response(w, 400, true, "Internal server error ("+err.Error()+")", map[string]any{})
		return
	}

	UserId := data.UserId
	JobId := data.JobId

	if UserId == "" {
		helper.Logger("error", "In Server: user_id is required")
		helper.Response(w, 400, true, "user_id is required", map[string]any{})
		return
	}

	if JobId == "" {
		helper.Logger("error", "In Server: job_id is required")
		helper.Response(w, 400, true, "job_id is required", map[string]any{})
		return
	}

	if UserId == "" {
		helper.Logger("error", "In Server: user_id is required")
		helper.Response(w, 400, true, "user_id is required", map[string]any{})
		return
	}

	data.JobId = JobId
	data.UserId = UserId

	result, err := service.JobFavourite(data)

	if err != nil {
		helper.Response(w, 400, true, err.Error(), map[string]any{})
		return
	}

	helper.Logger("info", "Job Favourite success")
	helper.Response(w, http.StatusOK, false, "Successfully", result)
}
