package response

import (
	"encoding/json"
	"net/http"
)

// Response represents the JSON response structure
type Response struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data,omitempty"`
	Error  string      `json:"error,omitempty"`
}

// ConstructResponse constructs and sends a JSON response
func ConstructResponse(w http.ResponseWriter, status int, data interface{}, err error) {
	response := Response{}

	if err != nil {
		response.Status = "error"
		response.Error = err.Error()
	} else {
		response.Data = data
		response.Status = "success"
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(response)
}
