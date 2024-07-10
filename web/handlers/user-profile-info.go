package handlers

import (
	"fmt"
	"librarymanagement/db"
	"librarymanagement/logger"
	"librarymanagement/web/utils"
	"log/slog"
	"net/http"
)

func UserInfo(w http.ResponseWriter, r *http.Request) {

	userId, err := getUserId(r)
	if err != nil {
		slog.Error("Failed to get user id", logger.Extra(map[string]any{
			"error":   err.Error(),
			"payload": userId,
		}))
		utils.SendError(w, http.StatusBadRequest, err.Error())
		return
	}

	var user *db.Reader
	if user, err = db.GetReaderRepo().GetUserInfo(*userId); err != nil {
		utils.SendError(w, http.StatusInternalServerError, err.Error())
		return
	}
	fmt.Println(user)
	utils.SendData(w, map[string]interface{}{
		"username": user.Name,
		"email":    user.Email,
	})
}
