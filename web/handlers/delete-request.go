package handlers

import (
	"librarymanagement/db"
	"librarymanagement/logger"
	"librarymanagement/web/utils"
	"log/slog"
	"net/http"
)

func DeleteBorrowRequest(w http.ResponseWriter, r *http.Request) {

	reqId, err := convIntToString(r.URL.Query().Get("request_id"))
	if err != nil {
		slog.Error("can not found borrow request_id in param", logger.Extra(map[string]any{
			"error":   err.Error(),
			"payload": reqId,
		}))
		utils.SendError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := db.GetBookRequestRepo().DeleteBorrowRequest(reqId); err != nil {
		utils.SendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SendData(w, "Book Deleted")
}
