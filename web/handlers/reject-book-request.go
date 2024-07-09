package handlers

import (
	"encoding/json"
	"librarymanagement/db"
	"librarymanagement/logger"
	"librarymanagement/web/utils"
	"log/slog"
	"net/http"
)

type RequestId struct {
	Request_Id int `json:"request_id" validate:"required"`
}

func RejectRequest(w http.ResponseWriter, r *http.Request) {

	var reqId RequestId
	if err := json.NewDecoder(r.Body).Decode(&reqId); err != nil {
		slog.Error("Failed to decode request id", logger.Extra(map[string]any{
			"error":   err.Error(),
			"payload": reqId,
		}))
		utils.SendError(w, http.StatusPreconditionFailed, err.Error())
		return
	}

	if err := utils.ValidateStruct(reqId); err != nil {
		slog.Error("Failed to validate new book data", logger.Extra(map[string]any{
			"error":   err.Error(),
			"payload": reqId,
		}))
		utils.SendError(w, http.StatusExpectationFailed, err.Error())
		return
	}

	if err := db.GetBookRequestRepo().ChangeBorrowStatus(reqId.Request_Id, "Reject"); err != nil {
		utils.SendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SendData(w, "Reject Request Executed Successfully")
}
