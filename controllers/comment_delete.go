package controllers

import (
	"encoding/json"
	"net/http"
	helper "superapps/helpers"
	models "superapps/models"
	service "superapps/services"
)

func CommentDelete(w http.ResponseWriter, r *http.Request) {
	data := &models.CommentDelete{}

	err := json.NewDecoder(r.Body).Decode(data)

	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		helper.Response(w, 400, true, "Internal server error ("+err.Error()+")", map[string]any{})
		return
	}

	Id := data.Id

	data.Id = Id

	result, err := service.CommentDelete(data)

	if err != nil {
		helper.Response(w, 400, true, err.Error(), map[string]interface{}{})
		return
	}

	helper.Logger("info", "Comment Delete success")
	helper.Response(w, http.StatusOK, false, "Successfully", result)
}
