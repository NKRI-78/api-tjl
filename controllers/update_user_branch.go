package controllers

import (
	"encoding/json"
	"net/http"
	"superapps/entities"
	helper "superapps/helpers"
	service "superapps/services"
)

func UpdateUserBranch(w http.ResponseWriter, r *http.Request) {
	data := &entities.UpdateUserBranch{}

	err := json.NewDecoder(r.Body).Decode(data)

	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		helper.Response(w, 400, true, "Internal server error ("+err.Error()+")", map[string]any{})
		return
	}

	fullname := data.Fullname
	phone := data.Phone
	password := data.Password
	roleId := data.RoleId

	if fullname == "" {
		helper.Logger("error", "In Server: fullname is required")
		helper.Response(w, 400, true, "fullname is required", map[string]any{})
		return
	}

	if phone == "" {
		helper.Logger("error", "In Server: phone is required")
		helper.Response(w, 400, true, "phone is required", map[string]any{})
		return
	}

	if password == "" {
		helper.Logger("error", "In Server: password is required")
		helper.Response(w, 400, true, "password is required", map[string]any{})
		return
	}

	if roleId == "" {
		helper.Logger("error", "In Server: role_id is required")
		helper.Response(w, 400, true, "role_id is required", map[string]any{})
		return
	}

	result, err := service.UpdateUserBranch(data)

	if err != nil {
		helper.Response(w, 400, true, err.Error(), map[string]any{})
		return
	}

	helper.Logger("info", "Register User Branch success")
	helper.Response(w, http.StatusOK, false, "Successfully", result)
}
