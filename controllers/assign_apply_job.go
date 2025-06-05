package controllers

import (
	"encoding/json"
	"net/http"
	"superapps/entities"
	helper "superapps/helpers"
	"superapps/services"
)

func AssignApplyJob(w http.ResponseWriter, r *http.Request) {

	data := &entities.AssignApplyJob{}

	err := json.NewDecoder(r.Body).Decode(data)

	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		helper.Response(w, 400, true, "Internal server error ("+err.Error()+")", map[string]any{})
		return
	}

	JobId := data.JobId
	UserId := data.UserId

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

	result, err := services.AssignApplyJob(data)

	if err != nil {
		helper.Response(w, 400, true, err.Error(), map[string]any{})
		return
	}

	helper.Logger("info", "Assign Apply Job success")
	helper.Response(w, http.StatusOK, false, "Successfully", result["data"])
}
