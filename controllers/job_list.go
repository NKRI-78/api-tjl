package controllers

import (
	"net/http"
	helper "superapps/helpers"
	"superapps/services"
)

func JobList(w http.ResponseWriter, r *http.Request) {
	page := r.URL.Query().Get("page")
	limit := r.URL.Query().Get("limit")

	salary := r.URL.Query().Get("salary")
	country := r.URL.Query().Get("country")
	position := r.URL.Query().Get("position")
	search := r.URL.Query().Get("search")
	userId := r.URL.Query().Get("user_id")

	result, err := services.JobList(userId, search, salary, country, position, page, limit)

	if err != nil {
		helper.Response(w, 400, true, err.Error(), map[string]any{})
		return
	}

	helper.Logger("info", "Job List success")
	helper.ResponseWithPagination(w, http.StatusOK, false, "Successfully",
		result["total"],
		result["per_page"],
		result["prev_page"],
		result["next_page"],
		result["current_page"],
		result["next_url"],
		result["prev_url"],
		result["data"],
	)
}
