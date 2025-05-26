package controllers

import (
	"encoding/json"
	"net/http"
	"superapps/entities"
	helper "superapps/helpers"
	"superapps/services"

	uuid "github.com/satori/go.uuid"
)

func ImportUserCreate(w http.ResponseWriter, r *http.Request) {

	data := &entities.ImportUserStore{}

	err := json.NewDecoder(r.Body).Decode(data)

	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		helper.Response(w, 400, true, "Internal server error ("+err.Error()+")", map[string]any{})
		return
	}

	data.UserId = uuid.NewV4().String()

	hashedPassword, err := helper.Hash(data.Password)
	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
	}

	data.Password = string(hashedPassword)

	Email := data.Email
	Password := data.Password
	Phone := data.Phone

	if Email == "" {
		helper.Logger("error", "In Server: email is required")
		helper.Response(w, 400, true, "email is required", map[string]any{})
		return
	}

	if Password == "" {
		helper.Logger("error", "In Server: password is required")
		helper.Response(w, 400, true, "password is required", map[string]any{})
		return
	}

	if Phone == "" {
		helper.Logger("error", "In Server: phone is required")
		helper.Response(w, 400, true, "phone is required", map[string]any{})
		return
	}

	result, err := services.ImportUserCreate(data)

	if err != nil {
		helper.Response(w, 400, true, err.Error(), map[string]any{})
		return
	}

	helper.Logger("info", "Store Icon success")
	helper.Response(w, http.StatusOK, false, "Successfully", result)
}
