package utils

import "net/http"

type Page struct {
	Items        interface{} `json:"items"`
	ItemsPerPage int         `json:"itemsPerPage"`
	PageNumber   int         `json:"pageNumber"`
	TotalItems   int         `json:"totalItems"`
	TotalPages   int         `json:"totalPages"`
}

func SendPage(w http.ResponseWriter, page Page) {
	SendData(w, page)
}
