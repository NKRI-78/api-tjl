package controllers

import (
	"net/http"
	helper "superapps/helpers"
	"superapps/services"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
)

func GetDocumentAdditional(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	typeParam := vars["type"]

	tokenHeader := r.Header.Get("Authorization")

	token := helper.DecodeJwt(tokenHeader)

	claims, _ := token.Claims.(jwt.MapClaims)

	userId, _ := claims["id"].(string)

	result, err := services.GetDocumentAdditional(userId, typeParam)

	if err != nil {
		helper.Response(w, 400, true, err.Error(), map[string]any{})
		return
	}

	helper.Logger("info", "Get Document Additional success")
	helper.Response(w, http.StatusOK, false, "Successfully",
		result["data"],
	)
}
