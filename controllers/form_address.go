package controllers

import (
	"encoding/json"
	"net/http"
	helper "superapps/helpers"
	"superapps/models"
	"superapps/services"

	"github.com/dgrijalva/jwt-go"
)

func FormAddress(w http.ResponseWriter, r *http.Request) {

	data := &models.FormPlace{}

	err := json.NewDecoder(r.Body).Decode(data)

	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		helper.Response(w, 400, true, "Internal server error ("+err.Error()+")", map[string]interface{}{})
		return
	}

	tokenHeader := r.Header.Get("Authorization")

	token := helper.DecodeJwt(tokenHeader)

	claims, _ := token.Claims.(jwt.MapClaims)

	userId, _ := claims["id"].(string)

	Province := data.ProvinceId
	City := data.CityId
	District := data.DistrictId
	Subdistrict := data.SubdistrictId
	DetailAddress := data.DetailAddress

	UserId := userId

	if Province == "" {
		helper.Logger("error", "In Server: province is required")
		helper.Response(w, 400, true, "province is required", map[string]interface{}{})
		return
	}

	if City == "" {
		helper.Logger("error", "In Server: city is required")
		helper.Response(w, 400, true, "city is required", map[string]interface{}{})
		return
	}

	if District == "" {
		helper.Logger("error", "In Server: district is required")
		helper.Response(w, 400, true, "district is required", map[string]interface{}{})
		return
	}

	if Subdistrict == "" {
		helper.Logger("error", "In Server: subdistrict is required")
		helper.Response(w, 400, true, "subdistrict is required", map[string]any{})
		return
	}

	if DetailAddress == "" {
		helper.Logger("error", "In Server: detail_address is required")
		helper.Response(w, 400, true, "detail_address is required", map[string]interface{}{})
		return
	}

	data.UserId = UserId

	result, err := services.FormPlace(data)

	if err != nil {
		helper.Response(w, 400, true, err.Error(), map[string]interface{}{})
		return
	}

	helper.Logger("info", "Form Region success")
	helper.Response(w, http.StatusOK, false, "Successfully", result)
}
