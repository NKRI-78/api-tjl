package controllers

import (
	"encoding/json"
	"net/http"
	"superapps/entities"
	helper "superapps/helpers"
	"superapps/services"
)

func EventStoreImage(w http.ResponseWriter, r *http.Request) {

	data := &entities.EventStoreImage{}

	err := json.NewDecoder(r.Body).Decode(data)

	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		helper.Response(w, 400, true, "Internal server error ("+err.Error()+")", map[string]any{})
		return
	}

	EventId := data.EventId
	Path := data.Path

	if EventId == "" {
		helper.Logger("error", "In Server: event_id is required")
		helper.Response(w, 400, true, "event_id is required", map[string]any{})
		return
	}

	if Path == "" {
		helper.Logger("error", "In Server: path is required")
		helper.Response(w, 400, true, "path is required", map[string]any{})
		return
	}

	data.EventId = EventId
	data.Path = Path

	result, err := services.EventStoreImage(data)

	if err != nil {
		helper.Response(w, 400, true, err.Error(), map[string]any{})
		return
	}

	helper.Logger("info", "Event Store Image success")
	helper.Response(w, http.StatusOK, false, "Successfully", result)
}
