package controllers

import (
	"encoding/json"
	"net/http"
	helper "superapps/helpers"
	"superapps/models"
	service "superapps/services"
)

func Register(w http.ResponseWriter, r *http.Request) {
	data := &models.User{}

	err := json.NewDecoder(r.Body).Decode(data)

	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		helper.Response(w, 400, true, "Internal server error ("+err.Error()+")", map[string]any{})
		return
	}

	jobId := data.JobId
	branchId := data.BranchId
	avatar := data.Avatar
	fullname := data.Fullname
	email := data.Email
	password := data.Password
	phone := data.Phone

	if jobId == "" {
		helper.Logger("error", "In Server: job_id is required")
		helper.Response(w, 400, true, "job_id is required", map[string]any{})
		return
	}

	if branchId == "" {
		helper.Logger("error", "In Server: branch_id is required")
		helper.Response(w, 400, true, "branch_id is required", map[string]any{})
		return
	}

	if avatar == "" {
		helper.Logger("error", "In Server: avatar is required")
		helper.Response(w, 400, true, "avatar is required", map[string]any{})
		return
	}

	if fullname == "" {
		helper.Logger("error", "In Server: fullname is required")
		helper.Response(w, 400, true, "fullname field is required", map[string]any{})
		return
	}

	if email == "" {
		helper.Logger("error", "In Server: email field is required")
		helper.Response(w, 400, true, "email address field is required", map[string]any{})
		return
	}

	validateEmail := helper.IsValidEmail(email)

	if !validateEmail {
		helper.Logger("error", "In Server: E-mail address is invalid")
		helper.Response(w, 400, true, "email address is invalid", map[string]any{})
		return
	}

	if phone == "" {
		helper.Logger("error", "In Server: phone field is required")
		helper.Response(w, 400, true, "phone field is required", map[string]any{})
		return
	}

	if password == "" {
		helper.Logger("error", "In Server: password field is required")
		helper.Response(w, 400, true, "password field is required", map[string]any{})
		return
	}

	result, err := service.Register(data)

	if err != nil {
		helper.Response(w, 400, true, err.Error(), map[string]any{})
		return
	}

	helper.Logger("info", "Register success")
	helper.Response(w, http.StatusOK, false, "Successfully", result)
}
