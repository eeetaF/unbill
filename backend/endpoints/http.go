package endpoints

import (
	"encoding/json"
	"log"
	"net/http"
)

// JSONResponse defines a consistent JSON structure
type JSONResponse struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
}

// WriteJSON sends a JSON-formatted response
func WriteJSON(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(JSONResponse{
		Message: message,
		Status:  statusCode,
	})
	if err != nil {
		return
	}
	writeToLog(statusCode, message)
}

func OKnocontent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}

func OK(w http.ResponseWriter, content interface{}) {
	jsonRender(w, http.StatusOK, content)
}

func writeToLog(status int, message string) {
	log.Printf("%v: %v\n", status, message)
}

// jsonRender render content as jsonRender
func jsonRender(w http.ResponseWriter, status int, content interface{}) {
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.WriteHeader(status)
	if content == nil {
		plainText(w, status, "{}")
		return
	}
	if err := json.NewEncoder(w).Encode(content); err != nil {
		writeToLog(500, err.Error())
	}
}

// plainText renders content as plain text
func plainText(w http.ResponseWriter, status int, text string) {
	plainTextInternal(w, status, []byte(text))
}

func plainTextInternal(w http.ResponseWriter, status int, bs []byte) {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "text/plain;charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	if _, err := w.Write(bs); err != nil {
		writeToLog(500, err.Error())
	}
}
