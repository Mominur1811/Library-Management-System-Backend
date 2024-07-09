package handlers

import (
	"librarymanagement/db"
	"librarymanagement/web/middlewire"
	"librarymanagement/web/utils"
	"net/http"
)



func UserInfo(w http.ResponseWriter, r *http.Request) {

	userId, err := middlewire.GetUserId(r)
	if err != nil {
		utils.SendError(w, http.StatusExpectationFailed, err.Error())
		return
	}

	var user *db.Reader
	if user, err = db.GetReaderRepo().GetUserInfo(*userId); err != nil {
		utils.SendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SendData(w, map[string]interface{}{
		"username": user.Name,
		"email": user.Email,
	})
}
