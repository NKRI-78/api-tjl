package controllers 

import (
	"net/http"
	service "superapps/services"
	helper "superapps/helpers"
)

func GetMembernear(w http.ResponseWriter, r *http.Request) {

	appName := r.Header.Get("APP_NAME")

	originLat  := r.URL.Query().Get("lat")
	originLng  := r.URL.Query().Get("lng")

	if appName == "" {
		helper.Logger("error", "In Server: APP_NAME headers is required")
		helper.Response(w, 400, true, "APP_NAME headers is required", map[string]interface{}{})
		return
	}

	if originLat == "" {
		helper.Logger("error", "In Server: lat is required")
		helper.Response(w, 400, true, "lat is required", map[string]interface{}{})
		return
	}

	if originLng == "" {
		helper.Logger("error", "In Server: lng is required")
		helper.Response(w, 400, true, "lng is required", map[string]interface{}{})
		return
	}

	result, err := service.GetMembernear(originLat, originLng, appName)

	if err != nil {
		helper.Response(w, 400, true, err.Error(), map[string]interface{}{})
		return
	}

	helper.Response(w, http.StatusOK, false, "Successfully", result["data"])
}	