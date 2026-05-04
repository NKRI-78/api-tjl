package controllers

import (
	"net/http"
	helper "superapps/helpers"
	"superapps/services"
)

func DetectFullbody(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(20 << 20) // 20MB
	if err != nil {
		helper.Response(w, 400, true, "Invalid form data", map[string]any{})
		return
	}

	folder := r.FormValue("folder")
	if folder == "" {
		helper.Response(w, 400, true, "Field folder is required", map[string]any{})
		return
	}

	subfolder := r.FormValue("subfolder")
	if folder == "" {
		helper.Response(w, 400, true, "Field folder is required", map[string]any{})
		return
	}

	file, fileHeader, err := r.FormFile("media")
	if err != nil {
		helper.Response(w, 400, true, "Field media is required", map[string]any{})
		return
	}
	defer file.Close()

	result, err := services.DetectFullbody(folder, subfolder, file, fileHeader)
	if err != nil {
		helper.Response(w, 400, true, err.Error(), map[string]any{})
		return
	}

	helper.Logger("info", "Detect fullbody success")
	helper.Response(w, http.StatusOK, false, result.Message, result.Data)
}
