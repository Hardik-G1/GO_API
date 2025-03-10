package sendresponse

import (
	"encoding/json"
	"net/http"
)

func SendResponse(w http.ResponseWriter, err bool, code string) {
	var ResponseTag = map[string]interface{}{
		"error": false,
		"code":  "",
	}
	w.Header().Set("Content-Type", "application/json")
	ResponseTag["error"] = err
	ResponseTag["code"] = code
	json.NewEncoder(w).Encode(ResponseTag)
	return
}
