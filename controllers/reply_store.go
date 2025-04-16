package controllers

import (
	"encoding/json"
	"net/http"
	"superapps/entities"
	helper "superapps/helpers"
	service "superapps/services"

	"github.com/dgrijalva/jwt-go"
	uuid "github.com/satori/go.uuid"
)

func ReplyStore(w http.ResponseWriter, r *http.Request) {
	data := &entities.ReplyStore{}

	err := json.NewDecoder(r.Body).Decode(data)

	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		helper.Response(w, 400, true, "Internal server error ("+err.Error()+")", map[string]any{})
		return
	}

	tokenHeader := r.Header.Get("Authorization")

	token := helper.DecodeJwt(tokenHeader)

	claims, _ := token.Claims.(jwt.MapClaims)

	userId, _ := claims["id"].(string)

	Id := uuid.NewV4().String()
	CommentId := data.CommentId
	Reply := data.Reply

	data.Id = Id
	data.UserId = userId

	if CommentId == "" {
		helper.Logger("error", "In Server: comment_id is required")
		helper.Response(w, 400, true, "comment_id is required", map[string]any{})
		return
	}

	if Reply == "" {
		helper.Logger("error", "In Server: reply is required")
		helper.Response(w, 400, true, "reply is required", map[string]any{})
		return
	}

	result, err := service.ReplyStore(data)

	if err != nil {
		helper.Response(w, 400, true, err.Error(), map[string]any{})
		return
	}

	helper.Logger("info", "Reply store success")
	helper.Response(w, http.StatusOK, false, "Successfully", result)
}
