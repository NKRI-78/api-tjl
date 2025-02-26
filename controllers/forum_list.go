package controllers

import (
	"net/http"
	helper "superapps/helpers"
	"superapps/services"
)

func ForumList(w http.ResponseWriter, r *http.Request) {
	page := r.URL.Query().Get("page")
	limit := r.URL.Query().Get("limit")
	search := r.URL.Query().Get("search")

	result, err := services.ForumList(search, page, limit)

	if err != nil {
		helper.Response(w, 400, true, err.Error(), map[string]interface{}{})
		return
	}

	println(result)

	helper.Logger("info", "Forum Delete success")
	helper.Response(w, http.StatusOK, false, "Successfully", nil)
}
