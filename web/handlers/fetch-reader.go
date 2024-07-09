package handlers

import (
	"librarymanagement/db"
	"librarymanagement/web/utils"
	"net/http"
)


func FetchUser(w http.ResponseWriter, r *http.Request) {

	readerList, err := db.GetReaderRepo().GetFetchUser()
	if err != nil {
		utils.SendError(w, http.StatusNotAcceptable, err.Error())
		return
	}

	utils.SendData(w, readerList)

}
