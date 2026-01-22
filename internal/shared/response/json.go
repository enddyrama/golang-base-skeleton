package response

import (
	"encoding/json"
	"net/http"
)

type APIResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Result  interface{} `json:"result"`
}

func JSON(w http.ResponseWriter, code int, message string, result interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	return json.NewEncoder(w).Encode(APIResponse{
		Code:    code,
		Message: message,
		Result:  result,
	})
}
