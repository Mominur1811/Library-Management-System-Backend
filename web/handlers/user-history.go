package handlers

import (
	"librarymanagement/db"
	"librarymanagement/logger"
	"librarymanagement/web/middlewire"
	"librarymanagement/web/utils"
	"log/slog"
	"net/http"
)

type HistoryFilterParam struct {
	SearchVal     string `json:"search"`
	RequestStatus string `json:"category"`
}

func UserHistory(w http.ResponseWriter, r *http.Request) {

	userId, err := getUserId(r)
	if err != nil {
		slog.Error("Failed to get user id", logger.Extra(map[string]any{
			"error":   err.Error(),
			"payload": userId,
		}))
		utils.SendError(w, http.StatusBadRequest, err.Error())
		return
	}

	historyFilterParam := utils.GetPaginationParams(r, defaultSortBy, defaultSortOrder)
	historyFilterParam.BorrowerId = *userId
	history, err := db.GetBookRequestRepo().GetBorrowStatus(historyFilterParam)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SendData(w, history)
}

func getUserId(r *http.Request) (*int, error) {

	var userId *int
	var err error
	userIdStr := r.URL.Query().Get("userid")

	if userIdStr == "" {
		userId, err = middlewire.GetUserId(r)
		if err != nil {
			return nil, err
		}
		return userId, err
	}

	var temp int
	temp, err = convIntToString(userIdStr)
	if err != nil {
		return nil, err
	}
	return &temp, err

}
