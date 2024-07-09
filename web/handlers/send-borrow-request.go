package handlers

import (
	"encoding/json"
	"librarymanagement/db"
	"librarymanagement/logger"
	"librarymanagement/web/middlewire"
	"librarymanagement/web/utils"
	"log/slog"
	"net/http"
)

func BorrowRequestBook(w http.ResponseWriter, r *http.Request) {

	// Fetch Request Book Data from Json
	var requestBook db.BorrowRequest
	var err error
	if err = json.NewDecoder(r.Body).Decode(&requestBook); err != nil {
		slog.Error("Failed to decode borrow request data", logger.Extra(map[string]any{
			"error":   err.Error(),
			"payload": requestBook,
		}))
		utils.SendError(w, http.StatusPreconditionFailed, err.Error())
		return
	}

	//Fetch Reader Id from Jwt Token
	readerId, err := middlewire.GetUserId(r)
	if err != nil {
		slog.Error("Failed to get User Id from jwt token", logger.Extra(map[string]any{
			"error":   err.Error(),
			"payload": readerId,
		}))
		utils.SendError(w, http.StatusPreconditionFailed, err.Error())
		return
	}

	//Assign Reader Id from Jwt Token
	requestBook.BorrowerId = *readerId

	//Validate data that has been fetched from json and jwt token
	if err = utils.ValidateStruct(requestBook); err != nil {
		slog.Error("Failed to validate request data", logger.Extra(map[string]any{
			"error":   err.Error(),
			"payload": requestBook,
		}))
		utils.SendError(w, http.StatusExpectationFailed, err.Error())
		return
	}

	//Validate Borrow Request. Check if that user already have or push request earlier.
	if err = db.GetBookRequestRepo().ValidateUserBorrowRequest(&requestBook); err != nil {
		utils.SendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	//Push Request in db.
	var pendingRequest *db.BorrowRequest
	if pendingRequest, err = db.GetBookRequestRepo().PushBorrowRequest(&requestBook); err != nil {
		utils.SendError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.SendData(w, pendingRequest)

}
