package handlers

import (
	"librarymanagement/db"
	"librarymanagement/web/utils"
	"net/http"
)

func DeleteAdmin(w http.ResponseWriter, r *http.Request) {

	email := r.URL.Query().Get("email")
	if email == "" {
		utils.SendError(w, http.StatusNotAcceptable, "Can not find Email in the url")
		return
	}

	if err := db.GetAdminRepo().DeleteAdmin(email); err != nil {
		utils.SendError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.SendData(w, "Deleted Successfully")
}
