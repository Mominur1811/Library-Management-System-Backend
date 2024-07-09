package handlers

import (
	"encoding/json"
	"librarymanagement/db"
	"librarymanagement/logger"
	"librarymanagement/web/utils"
	"log/slog"
	"net/http"
)

type ConfirmBookRequest struct {
	Request_Id int `json:"request_id" validate:"required"`
	Book_Id    int `json:"book_id"    validate:"required"`
}

func ApprovedBookRequest(w http.ResponseWriter, r *http.Request) {

	var reqId ConfirmBookRequest
	if err := json.NewDecoder(r.Body).Decode(&reqId); err != nil {
		slog.Error("Failed to decode Book Request Data (request_id, book_id)", logger.Extra(map[string]any{
			"error":   err.Error(),
			"payload": reqId,
		}))
		utils.SendError(w, http.StatusPreconditionFailed, err.Error())
		return
	}

	if err := utils.ValidateStruct(reqId); err != nil {
		slog.Error("Failed to validate book request data(request_id, book_id)", logger.Extra(map[string]any{
			"error":   err.Error(),
			"payload": reqId,
		}))
		utils.SendError(w, http.StatusExpectationFailed, err.Error())
		return
	}

	if err := db.GetBookRepo().UpdateBookAvailability(reqId.Book_Id); err != nil {
		utils.SendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := db.GetBookRequestRepo().ChangeBorrowStatus(reqId.Request_Id, "Approved"); err != nil {
		utils.SendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SendData(w, "Borrow Request Accepted Successfull")
}
