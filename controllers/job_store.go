package controllers

import (
	"encoding/json"
	"net/http"
	helper "superapps/helpers"
	"superapps/models"
	service "superapps/services"
)

func JobStore(w http.ResponseWriter, r *http.Request) {

	data := &models.JobStore{}

	err := json.NewDecoder(r.Body).Decode(data)

	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		helper.Response(w, 400, true, "Internal server error ("+err.Error()+")", map[string]any{})
		return
	}

	Title := data.Title
	Caption := data.Caption
	Salary := data.Salary
	CatId := data.CatId
	UserId := data.UserId

	// var IsDraft bool

	// if data.IsDraft == 1 {
	// 	IsDraft = true
	// } else {
	// 	IsDraft = false
	// }

	if Title == "" {
		helper.Logger("error", "In Server: title is required")
		helper.Response(w, 400, true, "title is required", map[string]any{})
		return
	}

	if Caption == "" {
		helper.Logger("error", "In Server: caption is required")
		helper.Response(w, 400, true, "caption is required", map[string]any{})
		return
	}

	if Salary == "" {
		helper.Logger("error", "In Server: salary is required")
		helper.Response(w, 400, true, "salary is required", map[string]any{})
		return
	}

	if CatId == "" {
		helper.Logger("error", "In Server: cat_id is required")
		helper.Response(w, 400, true, "cat_id is required", map[string]any{})
		return
	}

	if UserId == "" {
		helper.Logger("error", "In Server: user_id is required")
		helper.Response(w, 400, true, "user_id is required", map[string]any{})
		return
	}

	result, err := service.JobStore(data)

	if err != nil {
		helper.Response(w, 400, true, err.Error(), map[string]any{})
		return
	}

	helper.Logger("info", "Store Bannner success")
	helper.Response(w, http.StatusOK, false, "Successfully", result)
}
