package handlers

import (
	"encoding/json"
	"librarymanagement/config"
	"librarymanagement/db"
	"librarymanagement/logger"
	"librarymanagement/web/middlewire"
	"librarymanagement/web/utils"
	"log/slog"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func LoginReader(w http.ResponseWriter, r *http.Request) {

	var loginReader db.LoginReaderCredintials
	if err := json.NewDecoder(r.Body).Decode(&loginReader); err != nil {
		slog.Error("Failed to decode login data", logger.Extra(map[string]any{
			"error":   err.Error(),
			"payload": loginReader,
		}))
		utils.SendError(w, http.StatusPreconditionFailed, err.Error())
		return
	}

	if err := utils.ValidateStruct(loginReader); err != nil {
		slog.Error("Failed to validate login data", logger.Extra(map[string]any{
			"error":   err.Error(),
			"payload": loginReader,
		}))
		utils.SendError(w, http.StatusExpectationFailed, err.Error())
		return
	}

	loginReader.Password = hashPassword(loginReader.Password)

	var reader *db.Reader
	var err error
	if reader, err = db.GetReaderRepo().ValidateUserLogin(loginReader); err != nil {
		utils.SendError(w, http.StatusPreconditionFailed, err.Error())
		return
	}

	jwtToken, err := createToken(*reader.Id, 50, "user")
	if err != nil {
		slog.Error("Failed to get access token", logger.Extra(map[string]any{
			"error":     err.Error(),
			"jwt_token": jwtToken,
		}))
		utils.SendError(w, http.StatusExpectationFailed, err.Error())
		return
	}
	utils.SendData(w, map[string]interface{}{
		"username":  reader.Name,
		"email":     reader.Email,
		"jwt_token": jwtToken,
		"role":"user",
	})

}

// JWT CREATION
func createToken(userId int, lifeTime int, role string) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		middlewire.AuthClaims{
			Id:   userId,
			Role: role,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * time.Duration(lifeTime))),
			},
		})
	tokenString, err := token.SignedString([]byte(config.GetConfig().JwtSecretKey))
	if err != nil {
		return "", err
	}

	return "bearer " + tokenString, nil
}
