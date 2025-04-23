package controllers

import (
	"net/http"
	helper "superapps/helpers"
)

func SendEmail(w http.ResponseWriter, r *http.Request) {

	helper.SendEmail("reihanagam7@gmail.com", "", "", "")

	helper.Response(w, http.StatusOK, false, "Successfully", map[string]any{})
}
