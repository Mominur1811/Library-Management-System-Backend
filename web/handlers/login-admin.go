package handlers

import (
	"encoding/json"
	"librarymanagement/db"
	"librarymanagement/logger"
	"librarymanagement/web/utils"
	"log/slog"
	"net/http"
)

func LoginAdmin(w http.ResponseWriter, r *http.Request) {

	var loginAdmin db.LoginReaderCredintials
	if err := json.NewDecoder(r.Body).Decode(&loginAdmin); err != nil {
		slog.Error("Failed to decode login data", logger.Extra(map[string]any{
			"error":   err.Error(),
			"payload": LoginAdmin,
		}))
		utils.SendError(w, http.StatusPreconditionFailed, err.Error())
		return
	}

	if err := utils.ValidateStruct(loginAdmin); err != nil {
		slog.Error("Failed to validate login data", logger.Extra(map[string]any{
			"error":   err.Error(),
			"payload": loginAdmin,
		}))
		utils.SendError(w, http.StatusExpectationFailed, err.Error())
		return
	}

	loginAdmin.Password = hashPassword(loginAdmin.Password)

	var role string
	var err error
	if role, err = db.GetAdminRepo().ValidateAdminLogin(loginAdmin); err != nil {
		utils.SendError(w, http.StatusPreconditionFailed, err.Error())
		return
	}

	jwtToken, err := createToken(1, 50, role)
	if err != nil {
		slog.Error("Failed to get access token", logger.Extra(map[string]any{
			"error":     err.Error(),
			"jwt_token": jwtToken,
		}))
		utils.SendError(w, http.StatusExpectationFailed, err.Error())
		return
	}
	utils.SendData(w, map[string]interface{}{
		"username":  "Admin",
		"email":     loginAdmin.Email,
		"jwt_token": jwtToken,
		"role":      role,
	})

}
