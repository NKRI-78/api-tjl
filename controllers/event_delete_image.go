package controllers

import (
	"encoding/json"
	"net/http"
	"superapps/entities"
	helper "superapps/helpers"
	service "superapps/services"
)

func EventDeleteImage(w http.ResponseWriter, r *http.Request) {
	data := &entities.Event{}

	err := json.NewDecoder(r.Body).Decode(data)

	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		helper.Response(w, 400, true, "Internal server error ("+err.Error()+")", map[string]any{})
		return
	}

	Id := data.Id

	data.Id = Id

	result, err := service.EventDeleteImage(data)

	if err != nil {
		helper.Response(w, 400, true, err.Error(), map[string]any{})
		return
	}

	helper.Logger("info", "Event Delete Image success")
	helper.Response(w, http.StatusOK, false, "Successfully", result)
}
