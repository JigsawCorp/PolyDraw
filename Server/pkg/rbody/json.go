package rbody

import (
	"encoding/json"
	"net/http"
)

// JSON makes the response with payload as json format
func JSON(w http.ResponseWriter, status int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write([]byte(response))
}

// JSONError makes the error response with payload as json format
func JSONError(w http.ResponseWriter, code int, message string) {
	JSON(w, code, map[string]string{"error": message})
}
