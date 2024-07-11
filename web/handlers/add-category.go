package handlers

import (
	"encoding/json"
	"librarymanagement/db"
	"librarymanagement/logger"
	"librarymanagement/web/utils"
	"log/slog"
	"net/http"
)

func AddCategory(w http.ResponseWriter, r *http.Request) {
	var newCategory db.Category
	if err := json.NewDecoder(r.Body).Decode(&newCategory); err != nil {
		slog.Error("Failed to decode new user data", logger.Extra(map[string]any{
			"error":   err.Error(),
			"payload": newCategory,
		}))
		utils.SendError(w, http.StatusPreconditionFailed, err.Error())
		return
	}
	if err := utils.ValidateStruct(newCategory); err != nil {
		slog.Error("Failed to validate new book data", logger.Extra(map[string]any{
			"error":   err.Error(),
			"payload": newCategory,
		}))
		utils.SendError(w, http.StatusExpectationFailed, err.Error())
		return
	}

	var err error
	var insCategory *db.Category
	if insCategory, err = db.GetCategoryRepo().InsertCategory(&newCategory); err != nil {
		utils.SendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SendData(w, insCategory)
}
