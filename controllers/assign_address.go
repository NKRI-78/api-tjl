package controllers

import (
	"encoding/json"
	"net/http"
	"superapps/entities"
	helper "superapps/helpers"
	"superapps/services"
)

func AssignAddress(w http.ResponseWriter, r *http.Request) {

	data := &entities.AssignAddress{}

	err := json.NewDecoder(r.Body).Decode(data)

	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		helper.Response(w, 400, true, "Internal server error ("+err.Error()+")", map[string]any{})
		return
	}

	ProvinceId := data.ProvinceId
	CityId := data.CityId
	DistrictId := data.DistrictId
	Subdistrict := data.Subdistrict
	UserId := data.UserId

	if ProvinceId == "" {
		helper.Logger("error", "In Server: province_id is required")
		helper.Response(w, 400, true, "province_id is required", map[string]any{})
		return
	}

	if CityId == "" {
		helper.Logger("error", "In Server: city_id is required")
		helper.Response(w, 400, true, "city_id is required", map[string]any{})
		return
	}

	if DistrictId == "" {
		helper.Logger("error", "In Server: district_id is required")
		helper.Response(w, 400, true, "district_id is required", map[string]any{})
		return
	}

	if Subdistrict == "" {
		helper.Logger("error", "In Server: subdistrict_id is required")
		helper.Response(w, 400, true, "subdistrict_id is required", map[string]any{})
		return
	}

	if UserId == "" {
		helper.Logger("error", "In Server: user_id is required")
		helper.Response(w, 400, true, "user_id is required", map[string]any{})
		return
	}

	result, err := services.AssignAddress(data)

	if err != nil {
		helper.Response(w, 400, true, err.Error(), map[string]any{})
		return
	}

	helper.Logger("info", "Assign Address success")
	helper.Response(w, http.StatusOK, false, "Successfully", result["data"])
}
