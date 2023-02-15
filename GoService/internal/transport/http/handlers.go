package http

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func (rt *Router) getDataByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	uID := mux.Vars(r)["id"]

	data, err := rt.svc.GetDataByID(r.Context(), uID)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}
