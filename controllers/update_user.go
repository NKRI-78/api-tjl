package controllers

import (
	"encoding/json"
	"net/http"
	helper "superapps/helpers"
	"superapps/models"
	"superapps/services"
)

func UpdateUser(w http.ResponseWriter, r *http.Request) {

	data := &models.UpdateUser{}

	err := json.NewDecoder(r.Body).Decode(data)

	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		helper.Response(w, 400, true, "Internal server error ("+err.Error()+")", map[string]any{})
		return
	}

	UserId := data.UserId
	Email := data.Email
	Fullname := data.Fullname
	Gender := data.Gender
	Height := data.Height
	Weight := data.Weight
	Place := data.Place
	Religion := data.Religion
	MaritalStatus := data.MaritalStatus
	Birthdate := data.Birthdate

	if UserId == "" {
		helper.Logger("error", "In Server: user_id is required")
		helper.Response(w, 400, true, "user_id is required", map[string]any{})
		return
	}

	if Email == "" {
		helper.Logger("error", "In Server: email is required")
		helper.Response(w, 400, true, "email is required", map[string]any{})
		return
	}

	if Fullname == "" {
		helper.Logger("error", "In Server: fullname is required")
		helper.Response(w, 400, true, "fullname is required", map[string]any{})
		return
	}

	if Gender == "" {
		helper.Logger("error", "In Server: gender is required")
		helper.Response(w, 400, true, "gender is required", map[string]any{})
		return
	}

	if Height == "" {
		helper.Logger("error", "In Server: height is required")
		helper.Response(w, 400, true, "height is required", map[string]any{})
		return
	}

	if Weight == "" {
		helper.Logger("error", "In Server: weight is required")
		helper.Response(w, 400, true, "weight is required", map[string]any{})
		return
	}

	if Place == "" {
		helper.Logger("error", "In Server: place is required")
		helper.Response(w, 400, true, "place is required", map[string]any{})
		return
	}

	if Religion == "" {
		helper.Logger("error", "In Server: religion is required")
		helper.Response(w, 400, true, "religion is required", map[string]any{})
		return
	}

	if MaritalStatus == "" {
		helper.Logger("error", "In Server: marital_status is required")
		helper.Response(w, 400, true, "marital_status is required", map[string]any{})
		return
	}

	if Birthdate == "" {
		helper.Logger("error", "In Server: birthdate is required")
		helper.Response(w, 400, true, "birthdate is required", map[string]any{})
		return
	}

	result, err := services.UpdateUser(data)

	if err != nil {
		helper.Response(w, 400, true, err.Error(), map[string]any{})
		return
	}

	helper.Logger("info", "Update User success")
	helper.Response(w, http.StatusOK, false, "Successfully", result)
}
