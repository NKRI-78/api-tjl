package controllers

import (
	"encoding/json"
	"net/http"
	helper "superapps/helpers"
	"superapps/models"
	service "superapps/services"
)

func LoginAdmin(w http.ResponseWriter, r *http.Request) {

	data := &models.UserAdmin{}

	err := json.NewDecoder(r.Body).Decode(data)

	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		helper.Response(w, 400, true, "Internal server error ("+err.Error()+")", map[string]any{})
		return
	}

	val := data.Val
	password := data.Password

	if val == "" {
		helper.Logger("error", "In Server: val field is required")
		helper.Response(w, 400, true, "val field is required", map[string]any{})
		return
	}

	if password == "" {
		helper.Logger("error", "In Server: password field is required")
		helper.Response(w, 400, true, "password field is required", map[string]any{})
		return
	}

	result, err := service.LoginAdmin(data)

	if err != nil {
		helper.Response(w, 400, true, err.Error(), map[string]any{})
		return
	}

	helper.Logger("info", "Login Admin success")
	helper.Response(w, http.StatusOK, false, "Successfully", result)
}
