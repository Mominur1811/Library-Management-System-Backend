package handlers

import (
	"librarymanagement/db"
	"librarymanagement/web/utils"
	"net/http"
	"strconv"
)

func DeleteBook(w http.ResponseWriter, r *http.Request) {

	book_id := getBookId(r.URL.Query().Get("book_id"))
	if err := db.GetBookRepo().DeleteBook(book_id); err != nil {
		utils.SendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SendData(w, "Book Deleted")
}

func getBookId(str string) int {
	limit, _ := strconv.ParseInt(str, 10, 32)
	return int(limit)
}
