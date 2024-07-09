package handlers

import (
	"librarymanagement/db"
	"librarymanagement/web/utils"
	"net/http"
)

func FetchBorrowStatus(w http.ResponseWriter, r *http.Request) {

	//Fetch param from url
	fetchBorrowStatusParam := utils.GetPaginationParams(r, defaultSortBy, defaultSortOrder)
	borrowList, err := db.GetBookRequestRepo().GetBorrowStatus(fetchBorrowStatusParam)
	if err != nil {
		utils.SendError(w, http.StatusNotAcceptable, err.Error())
		return
	}

	utils.SendData(w, borrowList)
}
