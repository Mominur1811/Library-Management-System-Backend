package handlers

import (
	"librarymanagement/db"
	"librarymanagement/web/utils"
	"net/http"
)

type Email struct {
	Address string `json:"email"`
}

func GetUnapprovedUser(w http.ResponseWriter, r *http.Request) {

	readerList, err := db.GetReaderRepo().GetUnapprovedUser()
	if err != nil {
		utils.SendError(w, http.StatusNotAcceptable, err.Error())
		return
	}

	utils.SendData(w, readerList)

}
