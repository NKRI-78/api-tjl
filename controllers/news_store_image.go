package controllers

import (
	"encoding/json"
	"net/http"
	"superapps/entities"
	helper "superapps/helpers"
	"superapps/services"
)

func NewsStoreImage(w http.ResponseWriter, r *http.Request) {

	data := &entities.NewsStoreImage{}

	err := json.NewDecoder(r.Body).Decode(data)

	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		helper.Response(w, 400, true, "Internal server error ("+err.Error()+")", map[string]any{})
		return
	}

	NewsId := data.NewsId
	Path := data.Path

	if NewsId == "" {
		helper.Logger("error", "In Server: news_id is required")
		helper.Response(w, 400, true, "news_id is required", map[string]any{})
		return
	}

	if Path == "" {
		helper.Logger("error", "In Server: path is required")
		helper.Response(w, 400, true, "path is required", map[string]any{})
		return
	}

	data.NewsId = NewsId
	data.Path = Path

	result, err := services.NewsStoreImage(data)

	if err != nil {
		helper.Response(w, 400, true, err.Error(), map[string]any{})
		return
	}

	helper.Logger("info", "News Store success")
	helper.Response(w, http.StatusOK, false, "Successfully", result)
}
