package handlers

import (
	"encoding/json"
	"librarymanagement/db"
	"librarymanagement/logger"
	"librarymanagement/web/utils"
	"log/slog"
	"net/http"
)

func AddAdmin(w http.ResponseWriter, r *http.Request) {

	var newAdmin db.Admin
	if err := json.NewDecoder(r.Body).Decode(&newAdmin); err != nil {
		slog.Error("Failed to decode new user data", logger.Extra(map[string]any{
			"error":   err.Error(),
			"payload": newAdmin,
		}))
		utils.SendError(w, http.StatusPreconditionFailed, err.Error())
		return
	}

	if err := utils.ValidateStruct(newAdmin); err != nil {
		slog.Error("Failed to validate new user data", logger.Extra(map[string]any{
			"error":   err.Error(),
			"payload": newAdmin,
		}))
		utils.SendError(w, http.StatusExpectationFailed, err.Error())
		return
	}

	newAdmin.Password = hashPassword(newAdmin.Password)
	
	var insAdmin *db.Admin
	var err error
	if insAdmin, err = db.GetAdminRepo().RegisterAdmin(&newAdmin); err != nil {
		utils.SendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SendData(w, insAdmin)

}
