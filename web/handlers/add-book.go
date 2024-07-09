package handlers

import (
	"encoding/json"
	"librarymanagement/db"
	"librarymanagement/logger"
	"librarymanagement/web/utils"
	"log/slog"
	"net/http"
)

func AddBook(w http.ResponseWriter, r *http.Request) {
	var newBook db.Book
	if err := json.NewDecoder(r.Body).Decode(&newBook); err != nil {
		slog.Error("Failed to decode new user data", logger.Extra(map[string]any{
			"error":   err.Error(),
			"payload": newBook,
		}))
		utils.SendError(w, http.StatusPreconditionFailed, err.Error())
		return
	}
	if err := utils.ValidateStruct(newBook); err != nil {
		slog.Error("Failed to validate new book data", logger.Extra(map[string]any{
			"error":   err.Error(),
			"payload": newBook,
		}))
		utils.SendError(w, http.StatusExpectationFailed, err.Error())
		return
	}

	var insBook *db.Book
	var err error
	if insBook, err = db.GetBookRepo().InsertBook(&newBook); err != nil {
		utils.SendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SendData(w, insBook)
}
