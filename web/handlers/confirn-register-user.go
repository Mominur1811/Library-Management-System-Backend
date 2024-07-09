package handlers

import (
	"encoding/json"
	"librarymanagement/logger"
	"librarymanagement/web/utils"
	"log/slog"
	"net/http"
)

func AcceptRegistration(w http.ResponseWriter, r *http.Request) {

	var email string
	if err := json.NewDecoder(r.Body).Decode(&email); err != nil {
		slog.Error("Failed to decode email data", logger.Extra(map[string]any{
			"error":   err.Error(),
			"payload": email,
		}))
		utils.SendError(w, http.StatusPreconditionFailed, err.Error())
		return
	}

	//if err := db.GetAdminRepo().ConfirmRegistratio
}
