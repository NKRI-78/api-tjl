package controllers

import (
	"encoding/json"
	"net/http"
	"superapps/entities"
	helper "superapps/helpers"
	"superapps/services"
)

func AssignBranch(w http.ResponseWriter, r *http.Request) {

	data := &entities.AssignBranch{}

	err := json.NewDecoder(r.Body).Decode(data)

	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		helper.Response(w, 400, true, "Internal server error ("+err.Error()+")", map[string]any{})
		return
	}

	BranchId := data.BranchId
	UserId := data.UserId

	if BranchId == "" {
		helper.Logger("error", "In Server: branch_id is required")
		helper.Response(w, 400, true, "branch_id is required", map[string]any{})
		return
	}

	if UserId == "" {
		helper.Logger("error", "In Server: user_id is required")
		helper.Response(w, 400, true, "user_id is required", map[string]any{})
		return
	}

	result, err := services.AssignBranch(data)

	if err != nil {
		helper.Response(w, 400, true, err.Error(), map[string]any{})
		return
	}

	helper.Logger("info", "Assign Branch success")
	helper.Response(w, http.StatusOK, false, "Successfully", result["data"])
}
