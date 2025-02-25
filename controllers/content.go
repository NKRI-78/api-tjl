package controllers

import (
	"encoding/json"
	"net/http"
	helper "superapps/helpers"
	models "superapps/models"
	service "superapps/services"

	"github.com/dgrijalva/jwt-go"
)

func All(w http.ResponseWriter, r *http.Request) {

	page := r.URL.Query().Get("page")
	limit := r.URL.Query().Get("limit")
	search := r.URL.Query().Get("search")
	appName := r.URL.Query().Get("app_name")

	result, err := service.GetContent(search, page, limit, appName)

	if err != nil {
		helper.Response(w, 400, true, err.Error(), map[string]interface{}{})
		return
	}

	helper.Logger("info", "Get all content success")
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

func CreateMediaContent(w http.ResponseWriter, r *http.Request) {

	data := &models.ContentMedia{}

	errCreateMedia := json.NewDecoder(r.Body).Decode(data)

	if errCreateMedia != nil {
		helper.Logger("error", "In Server: "+errCreateMedia.Error())
		return
	}

	result, err := service.CreateMediaContent(data)

	if err != nil {
		helper.Response(w, 400, true, err.Error(), map[string]interface{}{})
		return
	}

	helper.Logger("info", "Create media success")
	helper.Response(w, http.StatusOK, false, "Successfully", result)
}

func CreateContent(w http.ResponseWriter, r *http.Request) {

	data := &models.Content{}

	errCreateContent := json.NewDecoder(r.Body).Decode(data)

	if errCreateContent != nil {
		helper.Logger("error", "In Server: "+errCreateContent.Error())
		return
	}

	tokenHeader := r.Header.Get("Authorization")

	token := helper.DecodeJwt(tokenHeader)

	claims, _ := token.Claims.(jwt.MapClaims)

	userId, _ := claims["id"].(string)

	appName := r.Header.Get("APP_NAME")
	id := data.Uid
	title := data.Title
	desc := data.Description

	TypeId := data.TypeId

	data.AppName = appName
	data.UserId = userId

	if appName == "" {
		helper.Logger("error", "In Server: APP_NAME headers is required")
		helper.Response(w, 400, true, "APP_NAME headers is required", map[string]interface{}{})
		return
	}

	if id == "" {
		helper.Logger("error", "In Server: id is required")
		helper.Response(w, 400, true, "id is required", map[string]interface{}{})
		return
	}

	if title == "" {
		helper.Logger("error", "In Server: title field is required")
		helper.Response(w, 400, true, "title field is required", map[string]interface{}{})
		return
	}

	if desc == "" {
		helper.Logger("error", "In Server: desc field is required")
		helper.Response(w, 400, true, "desc field is required", map[string]interface{}{})
		return
	}

	if TypeId == 0 {
		helper.Logger("error", "In Server: type_id field is required")
		helper.Response(w, 400, true, "type_id field is required", map[string]interface{}{})
		return
	}

	result, err := service.CreateContent(data)

	if err != nil {
		helper.Response(w, 400, true, err.Error(), map[string]interface{}{})
		return
	}

	helper.Logger("info", "Create content success")
	helper.Response(w, http.StatusOK, false, "Successfully", result)
}

func CreateContentLike(w http.ResponseWriter, r *http.Request) {

	data := &models.ReqContentLike{}

	errCreateContentLike := json.NewDecoder(r.Body).Decode(data)

	if errCreateContentLike != nil {
		helper.Logger("error", "In Server: "+errCreateContentLike.Error())
		return
	}

	tokenHeader := r.Header.Get("Authorization")

	token := helper.DecodeJwt(tokenHeader)

	claims, _ := token.Claims.(jwt.MapClaims)

	userId, _ := claims["id"].(string)

	contentId := data.ContentId

	data.ContentId = contentId
	data.UserId = userId

	if contentId == "" {
		helper.Logger("error", "In Server: content_id is required")
		helper.Response(w, 400, true, "content_id is required", map[string]interface{}{})
		return
	}

	_, err := service.CreateContentLike(data)

	if err != nil {
		helper.Response(w, 400, true, err.Error(), map[string]interface{}{})
		return
	}

	helper.Logger("info", "Create content like success")
	helper.Response(w, http.StatusOK, false, "Successfully", map[string]interface{}{})
}

func CreateContentUnlike(w http.ResponseWriter, r *http.Request) {

	data := &models.ReqContentUnlike{}

	errCreateContentUnlike := json.NewDecoder(r.Body).Decode(data)

	if errCreateContentUnlike != nil {
		helper.Logger("error", "In Server: "+errCreateContentUnlike.Error())
		return
	}

	tokenHeader := r.Header.Get("Authorization")

	token := helper.DecodeJwt(tokenHeader)

	claims, _ := token.Claims.(jwt.MapClaims)

	userId, _ := claims["id"].(string)

	contentId := data.ContentId

	data.ContentId = contentId
	data.UserId = userId

	if contentId == "" {
		helper.Logger("error", "In Server: content_id is required")
		helper.Response(w, 400, true, "content_id is required", map[string]interface{}{})
		return
	}

	_, err := service.CreateContentUnlike(data)

	if err != nil {
		helper.Response(w, 400, true, err.Error(), map[string]interface{}{})
		return
	}

	helper.Logger("info", "Create content unlike success")
	helper.Response(w, http.StatusOK, false, "Successfully", map[string]interface{}{})
}

func CreateContentComment(w http.ResponseWriter, r *http.Request) {

	data := &models.ReqContentComment{}

	errCreateContentComment := json.NewDecoder(r.Body).Decode(data)

	if errCreateContentComment != nil {
		helper.Logger("error", "In Server: "+errCreateContentComment.Error())
		return
	}

	tokenHeader := r.Header.Get("Authorization")

	token := helper.DecodeJwt(tokenHeader)

	claims, _ := token.Claims.(jwt.MapClaims)

	userId, _ := claims["id"].(string)

	contentId := data.ContentId
	comment := data.Comment

	data.ContentId = contentId
	data.Comment = comment
	data.UserId = userId

	if contentId == "" {
		helper.Logger("error", "In Server: content_id is required")
		helper.Response(w, 400, true, "content_id is required", map[string]interface{}{})
		return
	}

	if comment == "" {
		helper.Logger("error", "In Server: comment field is required")
		helper.Response(w, 400, true, "comment field is required", map[string]interface{}{})
		return
	}

	_, err := service.CreateContentComment(data)

	if err != nil {
		helper.Response(w, 400, true, err.Error(), map[string]interface{}{})
		return
	}

	helper.Logger("info", "Create content comment success")
	helper.Response(w, http.StatusOK, false, "Successfully", map[string]interface{}{})
}

func DeleteContent(w http.ResponseWriter, r *http.Request) {
	data := &models.DelContent{}

	errDelContent := json.NewDecoder(r.Body).Decode(data)

	if errDelContent != nil {
		helper.Logger("error", "In Server: "+errDelContent.Error())
		return
	}

	Uid := data.Uid

	data.Uid = Uid

	if Uid == "" {
		helper.Logger("error", "In Server: id field is required")
		helper.Response(w, 400, true, "id field is required", map[string]interface{}{})
		return
	}

	_, err := service.DelContent(data)

	if err != nil {
		helper.Response(w, 400, true, err.Error(), map[string]interface{}{})
		return
	}

	helper.Logger("info", "Delete content success")
	helper.Response(w, http.StatusOK, false, "Successfully", map[string]interface{}{})
}

func DeleteContentComment(w http.ResponseWriter, r *http.Request) {
	data := &models.DelContentComment{}

	errDelContentComment := json.NewDecoder(r.Body).Decode(data)

	if errDelContentComment != nil {
		helper.Logger("error", "In Server: "+errDelContentComment.Error())
		return
	}

	Uid := data.Uid

	data.Uid = Uid

	if Uid == "" {
		helper.Logger("error", "In Server: id field is required")
		helper.Response(w, 400, true, "id field is required", map[string]interface{}{})
		return
	}

	_, err := service.DelContentComment(data)

	if err != nil {
		helper.Response(w, 400, true, err.Error(), map[string]interface{}{})
		return
	}

	helper.Logger("info", "Delete content comment success")
	helper.Response(w, http.StatusOK, false, "Successfully", map[string]interface{}{})
}
