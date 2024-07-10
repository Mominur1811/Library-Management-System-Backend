package handlers

import (
	"encoding/json"
	"librarymanagement/db"
	"librarymanagement/logger"
	"librarymanagement/web/utils"
	"log/slog"
	"net/http"
)

func UpdateBook(w http.ResponseWriter, r *http.Request) {
	var updatedBookInfo db.Book
	if err := json.NewDecoder(r.Body).Decode(&updatedBookInfo); err != nil {
		slog.Error("Failed to decode new user data", logger.Extra(map[string]any{
			"error":   err.Error(),
			"payload": updatedBookInfo,
		}))
		utils.SendError(w, http.StatusPreconditionFailed, err.Error())
		return
	}

	if err := utils.ValidateStruct(updatedBookInfo); err != nil {
		slog.Error("Failed to validate new book data", logger.Extra(map[string]any{
			"error":   err.Error(),
			"payload": updatedBookInfo,
		}))
		utils.SendError(w, http.StatusExpectationFailed, err.Error())
		return
	}

	var err error
	updatedBookInfo.Id, err = convIntToString(r.PathValue("id"))
	if err != nil {
		slog.Error("can found book it in the path", logger.Extra(map[string]any{
			"error":   err.Error(),
			"payload": updatedBookInfo.Id,
		}))
		utils.SendError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err = db.GetBookRepo().UpdateBook(&updatedBookInfo); err != nil {
		utils.SendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SendData(w, "Updated")
}
