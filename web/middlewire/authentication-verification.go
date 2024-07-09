package middlewire

import (
	"fmt"
	"librarymanagement/config"
	"librarymanagement/logger"
	"librarymanagement/web/utils"
	"log/slog"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type UserId string

type AuthClaims struct {
	Id   int    `json:"Id"`
	Role string `json:"role"`
	jwt.RegisteredClaims
}

func unauthorizedResponse(w http.ResponseWriter) {
	utils.SendError(w, http.StatusUnauthorized, "Unauthorized")
}

func AuthenticateUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		str := r.Header.Get("authorization")
		tokenStr, err := ExtractToken(str)
		if err != nil {
			unauthorizedResponse(w)
			return
		}

		// parse jwt
		var claims AuthClaims
		token, err := jwt.ParseWithClaims(
			tokenStr,
			&claims,
			func(t *jwt.Token) (interface{}, error) {
				return []byte(config.GetConfig().JwtSecretKey), nil
			},
		)
		if err != nil || !token.Valid || claims.Role != "user" {
			unauthorizedResponse(w)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func AuthenticateAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		str := r.Header.Get("authorization")
		tokenStr, err := ExtractToken(str)
		if err != nil {
			slog.Error("Error is ExtractToken funtion", logger.Extra(map[string]any{
				"error":   err.Error(),
				"payload": tokenStr,
			}))
			unauthorizedResponse(w)
			return
		}

		// parse jwt
		var claims AuthClaims
		token, err := jwt.ParseWithClaims(
			tokenStr,
			&claims,
			func(t *jwt.Token) (interface{}, error) {
				return []byte(config.GetConfig().JwtSecretKey), nil
			},
		)
		if err != nil {
			slog.Error("Error generate authenticate admin", logger.Extra(map[string]any{
				"error":   err.Error(),
				"payload": token,
			}))
			unauthorizedResponse(w)
			return
		}

		if !token.Valid {
			slog.Error("Token is not valid", logger.Extra(map[string]any{
				"error":   fmt.Errorf("token is not valid"),
				"payload": token,
			}))
			unauthorizedResponse(w)
			return
		}

		if claims.Role != "super_admin" && claims.Role != "admin" {
			slog.Error("Role is not valid", logger.Extra(map[string]any{
				"error":   fmt.Errorf("invalid role"),
				"payload": claims,
			}))
			unauthorizedResponse(w)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func ExtractToken(header string) (string, error) {
	if len(header) == 0 {
		return "", fmt.Errorf("access token is null ")
	}

	//Check and Extract jwt part
	tokens := strings.Split(header, " ")
	if len(tokens) != 2 {
		return "", fmt.Errorf("access token structure is invalid ")
	}
	return tokens[1], nil

}

func GetUserId(r *http.Request) (*int, error) {

	str := r.Header.Get("authorization")
	tokenStr, err := ExtractToken(str)
	if err != nil {
		return nil, fmt.Errorf("error fetching jwt token")
	}
	var claims AuthClaims
	// parse jwt
	_, err = jwt.ParseWithClaims(
		tokenStr,
		&claims,
		func(t *jwt.Token) (interface{}, error) {
			return []byte(config.GetConfig().JwtSecretKey), nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("error claming info from token")
	}

	return &claims.Id, nil
}
