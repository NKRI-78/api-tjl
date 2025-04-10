package controllers

import (
	"encoding/json"
	"net/http"
	helper "superapps/helpers"
	"superapps/models"
	"superapps/services"
)

func BannerUpdate(w http.ResponseWriter, r *http.Request) {

	data := &models.Banner{}

	err := json.NewDecoder(r.Body).Decode(data)

	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		helper.Response(w, 400, true, "Internal server error ("+err.Error()+")", map[string]any{})
		return
	}

	Id := data.Id
	Path := data.Path
	Link := data.Link

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

	if Link == "" {
		helper.Logger("error", "In Server: link is required")
		helper.Response(w, 400, true, "link is required", map[string]any{})
		return
	}

	result, err := services.BannerUpdate(data)

	if err != nil {
		helper.Response(w, 400, true, err.Error(), map[string]any{})
		return
	}

	helper.Logger("info", "Store Bannner success")
	helper.Response(w, http.StatusOK, false, "Successfully", result)
}
