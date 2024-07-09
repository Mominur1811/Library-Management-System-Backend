package handlers

import (
	"encoding/json"
	"fmt"
	"librarymanagement/db"
	"librarymanagement/logger"
	"librarymanagement/web/utils"
	"log/slog"
	"net/http"
)

type ReadProgress struct {
	RequestId int `json:"request_id"`
	PageCount int `json:"page_cnt"`
}

func UserReadProgressUpdate(w http.ResponseWriter, r *http.Request) {

	var newProgress ReadProgress
	if err := json.NewDecoder(r.Body).Decode(&newProgress); err != nil {
		slog.Error("Failed to decode user todays read progress", logger.Extra(map[string]any{
			"error":   err.Error(),
			"payload": newProgress,
		}))
		utils.SendError(w, http.StatusPreconditionFailed, err.Error())
		return
	}

	//fmt.Println(newProgress)
    fmt.Println(newProgress)
	if err := db.GetBookRequestRepo().UpdateUserReadProgress(newProgress.RequestId, newProgress.PageCount); err != nil {
		utils.SendError(w, http.StatusNotAcceptable, err.Error())
		return
	}

	utils.SendData(w, "Successfully updated today's progress")
}
