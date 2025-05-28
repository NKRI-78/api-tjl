package controllers

import (
	"encoding/json"
	"net/http"
	helper "superapps/helpers"
	"superapps/models"

	"github.com/bregydoc/gtranslate"
)

func TranslateLang(w http.ResponseWriter, r *http.Request) {

	data := &models.TranslateRequest{}

	errBody := json.NewDecoder(r.Body).Decode(data)

	if errBody != nil {
		helper.Logger("error", "In Server: "+errBody.Error())
		helper.Response(w, 400, true, "Internal server error ("+errBody.Error()+")", map[string]any{})
		return
	}

	translated, errTranslate := gtranslate.TranslateWithParams(data.Text, gtranslate.TranslationParams{
		From: "auto",
		To:   data.To,
	})
	if errTranslate != nil {
		helper.Logger("error", "In Server: "+errTranslate.Error())
		helper.Response(w, 400, true, "Internal server error ("+errTranslate.Error()+")", map[string]any{})
		return
	}

	response := models.TranslateResponse{
		Text: translated,
	}

	helper.Logger("info", "Translate success")
	helper.Response(w, http.StatusOK, false, "Successfully", response)
}
