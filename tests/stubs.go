package tests

import (
	"encoding/json"
	"log"
	"net/http"
)

type MockResponseWriter struct {
}

func NewMockResponseWriter() http.ResponseWriter {
	return &MockResponseWriter{}
}

func (m MockResponseWriter) Header() http.Header {
	return map[string][]string{}
}

func (m MockResponseWriter) Write(_ []byte) (int, error) {
	return 0, nil
}

func (m MockResponseWriter) WriteHeader(_ int) {

}

func Handler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Welcome to yashpokar/router",
	}); err != nil {
		log.Fatal(err)
	}
}

func PanicMaker(_ http.ResponseWriter, _ *http.Request) {
	panic("to test panic recovery")
}
