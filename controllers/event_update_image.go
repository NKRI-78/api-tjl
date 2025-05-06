package controllers

import (
	"encoding/json"
	"net/http"
	"superapps/entities"
	helper "superapps/helpers"
	"superapps/services"
)

func EventUpdateImage(w http.ResponseWriter, r *http.Request) {

	data := &entities.EventUpdateImage{}

	err := json.NewDecoder(r.Body).Decode(data)

	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		helper.Response(w, 400, true, "Internal server error ("+err.Error()+")", map[string]any{})
		return
	}

	Id := data.Id
	Path := data.Path

	if Id == "" {
		helper.Logger("error", "In Server: id is required")
		helper.Response(w, 400, true, "id is required", map[string]any{})
		return
	}

	if Path == "" {
		helper.Logger("error", "In Server: path is required")
		helper.Response(w, 400, true, "path is required", map[string]any{})
		return
	}

	result, err := services.EventUpdateImage(data)

	if err != nil {
		helper.Response(w, 400, true, err.Error(), map[string]any{})
		return
	}

	helper.Logger("info", "Event Update Image success")
	helper.Response(w, http.StatusOK, false, "Successfully", result)
}
