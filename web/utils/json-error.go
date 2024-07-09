package utils

import "net/http"

func SendError(w http.ResponseWriter, status int, message string) {

	SendJson(w, status, map[string]interface{}{
		"status":  status,
		"message": message,
	})
}
