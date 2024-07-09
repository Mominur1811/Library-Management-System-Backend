package handlers

import (
	"librarymanagement/db"
	"librarymanagement/web/middlewire"
	"librarymanagement/web/utils"
	"net/http"
)

type HistoryFilterParam struct {
	SearchVal     string `json:"search"`
	RequestStatus string `json:"category"`
}

func UserHistory(w http.ResponseWriter, r *http.Request) {

	userId, err := middlewire.GetUserId(r)
	if err != nil {
		utils.SendError(w, http.StatusExpectationFailed, err.Error())
		return
	}

	historyFilterParam := utils.GetPaginationParams(r, defaultSortBy, defaultSortOrder)
	historyFilterParam.BorrowerId = *userId
	history, err := db.GetBookRequestRepo().GetBorrowStatus(historyFilterParam)
	if err != nil {
		utils.SendError(w, http.StatusNotAcceptable, err.Error())
		return
	}

	utils.SendData(w, history)
}
