package middlewares

import (
	"context"
	"net/http"
	"os"
	"strings"
	"github.com/dgrijalva/jwt-go"
	helper "superapps/helpers"
)

func JwtAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var url string = r.URL.Path

		var containImage string = "jpg"
	 
		if strings.Contains(url, containImage) {
		   next.ServeHTTP(w, r)
		   return 
		} 

		if r.URL.Path == "/api/v1/login" || r.URL.Path == "/api/v1/register" || r.URL.Path == "/api/v1/verify-otp" || r.URL.Path == "/api/v1/resend-otp" {
			next.ServeHTTP(w, r)
			return
		}

		tokenHeader := r.Header.Get("Authorization")

		if tokenHeader == "" {
			helper.Logger("error", "In Server: Missing auth token")
			helper.Response(w, http.StatusUnauthorized, true, "Missing auth token", map[string]interface{}{})
			return
		}

		splitted := strings.Split(tokenHeader, " ")

		if len(splitted) != 2 {
			helper.Logger("error", "In Server: Missing auth token")
			helper.Response(w, http.StatusUnauthorized, true, "Missing auth token", map[string]interface{}{})
			return
		}

		tokenPart := splitted[1]

		claims := jwt.MapClaims{}

		token, err := jwt.ParseWithClaims(tokenPart, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil {
			helper.Logger("error", "In Server: "+ err.Error())
			helper.Response(w, http.StatusUnauthorized, true, err.Error(), map[string]interface{}{})
			return
		}

		if !token.Valid {
			helper.Logger("error", "In Server: Token is invalid")
			helper.Response(w, http.StatusUnauthorized, true, "Token is invalid", map[string]interface{}{})
			return
		}

		ctx := context.WithValue(r.Context(), "user", claims)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func CreateToken(userId string) (map[string]string, error) {

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["authorized"] = true
	claims["id"] = userId
	// claims["exp"] = time.Now().Add(time.Hour * 168).Unix()

	access, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))

	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		return nil, err
	}

	return map[string]string{"token": access}, nil
}
