package api

import "net/http"

type healthResponse struct {
	Status  string `json:"status"`
	Service string `json:"service"`
}

func health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, healthResponse{
		Status:  "ok",
		Service: "graphsentinel",
	})
}
