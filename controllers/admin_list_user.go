package controllers

import (
	"net/http"
	helper "superapps/helpers"
	"superapps/services"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
)

func AdminListUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	Type := vars["type"]

	tokenHeader := r.Header.Get("Authorization")

	token := helper.DecodeJwt(tokenHeader)

	claims, _ := token.Claims.(jwt.MapClaims)

	BranchId, _ := claims["branch_id"].(string)

	result, err := services.AdminListUser(Type, BranchId)

	if err != nil {
		helper.Response(w, 400, true, err.Error(), map[string]any{})
		return
	}

	helper.Logger("info", "Get Admin List User success")
	helper.Response(w, http.StatusOK, false, "Successfully", result["data"])
}
