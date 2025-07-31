package controllers

import (
	"encoding/json"
	"net/http"
	"superapps/entities"
	helper "superapps/helpers"
	service "superapps/services"
)

func UpdateBranch(w http.ResponseWriter, r *http.Request) {
	data := &entities.UpdateBranch{}

	err := json.NewDecoder(r.Body).Decode(data)

	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		helper.Response(w, 400, true, "Internal server error ("+err.Error()+")", map[string]any{})
		return
	}

	BranchName := data.Name

	if BranchName == "" {
		helper.Logger("error", "In Server: name is required")
		helper.Response(w, 400, true, "name is required", map[string]any{})
		return
	}

	result, err := service.UpdateBranch(data)

	if err != nil {
		helper.Response(w, 400, true, err.Error(), map[string]any{})
		return
	}

	helper.Logger("info", "Update Branch success")
	helper.Response(w, http.StatusOK, false, "Successfully", result)
}
