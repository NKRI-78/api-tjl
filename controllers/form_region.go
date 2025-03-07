package controllers

import (
	"encoding/json"
	"net/http"
	helper "superapps/helpers"
	"superapps/models"
	"superapps/services"

	"github.com/dgrijalva/jwt-go"
)

func FormRegion(w http.ResponseWriter, r *http.Request) {

	data := &models.FormRegion{}

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

	ProvinceId := data.ProvinceId
	CityId := data.CityId
	DistrictId := data.DistrictId
	SubdistrictId := data.SubdistrictId
	DetailAddress := data.DetailAddress

	UserId := userId

	if ProvinceId == "" {
		helper.Logger("error", "In Server: province_id is required")
		helper.Response(w, 400, true, "province_id is required", map[string]interface{}{})
		return
	}

	if CityId == "" {
		helper.Logger("error", "In Server: city_id is required")
		helper.Response(w, 400, true, "city_id is required", map[string]interface{}{})
		return
	}

	if DistrictId == "" {
		helper.Logger("error", "In Server: district_id is required")
		helper.Response(w, 400, true, "district_id is required", map[string]interface{}{})
		return
	}

	if SubdistrictId == "" {
		helper.Logger("error", "In Server: height is required")
		helper.Response(w, 400, true, "height is required", map[string]any{})
		return
	}

	if DetailAddress == "" {
		helper.Logger("error", "In Server: detail_address is required")
		helper.Response(w, 400, true, "detail_address is required", map[string]interface{}{})
		return
	}

	data.UserId = UserId

	result, err := services.FormRegion(data)

	if err != nil {
		helper.Response(w, 400, true, err.Error(), map[string]interface{}{})
		return
	}

	helper.Logger("info", "Form Region success")
	helper.Response(w, http.StatusOK, false, "Successfully", result)
}
