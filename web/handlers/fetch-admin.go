package handlers

import (
	"librarymanagement/db"
	"librarymanagement/web/utils"
	"net/http"
)

func FetchAdmin(w http.ResponseWriter, r *http.Request) {

	adminList, err := db.GetAdminRepo().GetFetchAdmin()
	if err != nil {
		utils.SendError(w, http.StatusNotAcceptable, err.Error())
		return
	}
	utils.SendData(w, adminList)
}
