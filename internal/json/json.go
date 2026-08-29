package json

import ("net/http"
	"encoding/json"

)

func SendJson( w http.ResponseWriter, data any, statusCode int){
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}