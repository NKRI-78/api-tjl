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

	avatar := data.Avatar
	fullname := data.Fullname
	email := data.Email
	phone := data.Phone
	password := data.Password
	roleId := data.RoleId

	if avatar == "" {
		helper.Logger("error", "In Server: avatar is required")
		helper.Response(w, 400, true, "avatar is required", map[string]any{})
		return
	}

	if fullname == "" {
		helper.Logger("error", "In Server: fullname is required")
		helper.Response(w, 400, true, "fullname is required", map[string]any{})
		return
	}

	if email == "" {
		helper.Logger("error", "In Server: email is required")
		helper.Response(w, 400, true, "email is required", map[string]any{})
		return
	}

	if phone == "" {
		helper.Logger("error", "In Server: phone is required")
		helper.Response(w, 400, true, "phone is required", map[string]any{})
		return
	}

	validateEmail := helper.IsValidEmail(email)

	if !validateEmail {
		helper.Logger("error", "In Server: E-mail address is invalid")
		helper.Response(w, 400, true, "email address is invalid", map[string]any{})
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
