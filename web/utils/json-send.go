package utils

import (
	"encoding/json"
	"net/http"
)

func SendJson(w http.ResponseWriter, status int, data interface{}) {

	w.Header().Add("Content-Type", "application/json")

	str, err := json.Marshal(data)
	if err != nil {
		message := "Convertion to Json failed"
		SendError(w, status, message)
		return
	}

	w.WriteHeader(status)
	w.Write(str)

}
