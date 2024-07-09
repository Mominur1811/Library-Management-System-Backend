package handlers

import (
	"encoding/json"
	"librarymanagement/db"
	"librarymanagement/logger"
	"librarymanagement/web/utils"
	"log/slog"
	"net/http"
)

type Email struct {
	Address string `json:"email"`
}

func GetUnapprovedUser(w http.ResponseWriter, r *http.Request) {

	readerList, err := db.GetReaderRepo().GetUnapprovedUser()
	if err != nil {
		utils.SendError(w, http.StatusNotAcceptable, err.Error())
		return
	}

	utils.SendData(w,  readerList)

}

func ApprovedUser(w http.ResponseWriter, r *http.Request) {

	var eAdd Email
	if err := json.NewDecoder(r.Body).Decode(&eAdd); err != nil {
		slog.Error("Failed to decode new user data", logger.Extra(map[string]any{
			"error":   err.Error(),
			"payload": eAdd,
		}))
		utils.SendError(w, http.StatusPreconditionFailed, err.Error())
		return
	}

	if err := utils.ValidateStruct(eAdd); err != nil {
		slog.Error("Failed to validate new book data", logger.Extra(map[string]any{
			"error":   err.Error(),
			"payload": eAdd,
		}))
		utils.SendError(w, http.StatusExpectationFailed, err.Error())
		return
	}

	if err := db.GetReaderRepo().ApprovedReader(eAdd.Address); err != nil {
		utils.SendError(w, http.StatusPreconditionFailed, err.Error())
		return
	}

	utils.SendData(w, "Succeed")
}
