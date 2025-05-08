package controllers

import (
	"encoding/json"
	"net/http"
	"superapps/entities"
	helper "superapps/helpers"
	"superapps/services"
)

func ClientStore(w http.ResponseWriter, r *http.Request) {

	clientStore := &entities.ClientStoreResponse{}

	data := &entities.ClientStore{}

	err := json.NewDecoder(r.Body).Decode(data)

	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		helper.Response(w, 400, true, "Internal server error ("+err.Error()+")", map[string]any{})
		return
	}

	Icon := data.Icon
	Link := data.Link
	Name := data.Name

	if Icon == "" {
		helper.Logger("error", "In Server: icon is required")
		helper.Response(w, 400, true, "icon is required", map[string]any{})
		return
	}

	if Link == "" {
		helper.Logger("error", "In Server: link is required")
		helper.Response(w, 400, true, "link is required", map[string]any{})
		return
	}

	if Name == "" {
		helper.Logger("error", "In Server: name is required")
		helper.Response(w, 400, true, "name is required", map[string]any{})
		return
	}

	result, err := services.ClientStore(data)

	if err != nil {
		helper.Response(w, 400, true, err.Error(), map[string]any{})
		return
	}

	Id := result["data"].(int64)

	clientStore.Id = Id
	clientStore.Icon = Icon
	clientStore.Link = Link
	clientStore.Name = Name

	helper.Logger("info", "Client Store success")
	helper.Response(w, http.StatusOK, false, "Successfully", clientStore)
}
