package handlers

import (
	"encoding/json"
	"librarymanagement/db"
	"librarymanagement/logger"
	"librarymanagement/web/utils"
	"log/slog"
	"net/http"
)

func ReturnBook(w http.ResponseWriter, r *http.Request) {

	var reqId ConfirmBookRequest

	if err := json.NewDecoder(r.Body).Decode(&reqId); err != nil {
		slog.Error("Failed to decode user todays read progress", logger.Extra(map[string]any{
			"error":   err.Error(),
			"payload": reqId,
		}))
		utils.SendError(w, http.StatusPreconditionFailed, err.Error())
		return
	}
	if err := ExecuteReturnBook(reqId.Request_Id, reqId.Book_Id); err != nil {
		utils.SendError(w, http.StatusNotAcceptable, err.Error())
		return
	}

	utils.SendData(w, "Book Return Successfully!")
}

// Return Book Function First Update book_request table then Update book talbe
func ExecuteReturnBook(reqId int, bookId int) error {

	if err := db.GetBookRequestRepo().ChangeBorrowStatus(reqId, "Returned"); err != nil {
		return err
	}

	if err := db.GetBookRepo().UpdateBookCount(bookId); err != nil {
		return err
	}
	return nil
}
