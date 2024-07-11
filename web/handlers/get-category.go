package handlers

import (
	"librarymanagement/db"
	"librarymanagement/web/utils"
	"net/http"
)

func FetchCategory(w http.ResponseWriter, r *http.Request) {
	categorylist, err := db.GetCategoryRepo().GetCategory()
	if err != nil {
		utils.SendError(w, http.StatusFailedDependency, err.Error())
		return
	}

	utils.SendData(w, categorylist)

}
