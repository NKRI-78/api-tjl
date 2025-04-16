package controllers

import (
	"net/http"
	helper "superapps/helpers"
	"superapps/services"
)

func JobSkillCategoryList(w http.ResponseWriter, r *http.Request) {

	result, err := services.JobSkillCategoryList()

	if err != nil {
		helper.Response(w, 400, true, err.Error(), map[string]any{})
		return
	}

	helper.Logger("info", "Fetch Job Skill Category List success")
	helper.Response(w, http.StatusOK, false, "Successfully", result["data"])
}
