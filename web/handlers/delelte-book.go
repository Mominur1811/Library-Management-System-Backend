package handlers

import (
	"fmt"
	"librarymanagement/db"
	"librarymanagement/logger"
	"librarymanagement/web/utils"
	"log/slog"
	"net/http"
	"strconv"
)

func DeleteBook(w http.ResponseWriter, r *http.Request) {

	book_id, err := convIntToString(r.URL.Query().Get("book_id"))
	if err != nil {
		slog.Error("can not found book in param", logger.Extra(map[string]any{
			"error":   err.Error(),
			"payload": book_id,
		}))
		utils.SendError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := db.GetBookRepo().DeleteBook(book_id); err != nil {
		utils.SendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SendData(w, "Book Deleted")
}

func convIntToString(str string) (int, error) {
	if str == "" {
		return 0, fmt.Errorf("id not found")
	}
	limit, _ := strconv.ParseInt(str, 10, 32)
	return int(limit), nil
}
