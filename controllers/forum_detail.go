package controllers

import (
	"net/http"
	helper "superapps/helpers"
	models "superapps/models"
	"superapps/services"

	"github.com/gorilla/mux"
)

func ForumDetail(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	data := &models.Forum{}

	data.Id = id

	result, err := services.ForumDetail(data)

	if err != nil {
		helper.Response(w, 400, true, err.Error(), map[string]any{})
		return
	}

	helper.Logger("info", "Forum Detail success")
	helper.Response(w, http.StatusOK, false, "Successfully",
		result["data"],
	)
}
