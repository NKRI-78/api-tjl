package controllers

import (
	"encoding/json"
	"net/http"
	"superapps/entities"
	helper "superapps/helpers"
	"superapps/services"
)

func JobSkillCategoryStore(w http.ResponseWriter, r *http.Request) {

	data := &entities.JobSkillCategoryStore{}

	err := json.NewDecoder(r.Body).Decode(data)

	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		helper.Response(w, 400, true, "Internal server error ("+err.Error()+")", map[string]any{})
		return
	}

	CatId := data.CatId
	JobId := data.JobId

	if CatId == "" {
		helper.Logger("error", "In Server: cat_id is required")
		helper.Response(w, 400, true, "cat_id is required", map[string]any{})
		return
	}

	if JobId == "" {
		helper.Logger("error", "In Server: job_id is required")
		helper.Response(w, 400, true, "job_id is required", map[string]any{})
		return
	}

	result, err := services.JobSkillCategoryStore(data)

	if err != nil {
		helper.Response(w, 400, true, err.Error(), map[string]any{})
		return
	}

	helper.Logger("info", "Store Job Skill Category success")
	helper.Response(w, http.StatusOK, false, "Successfully", result["data"])
}
