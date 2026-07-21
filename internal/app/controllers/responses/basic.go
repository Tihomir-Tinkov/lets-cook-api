package responses

import (
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog/log"
)

func JSONError(w http.ResponseWriter, r *http.Request, err error, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	log.Error().Err(err).Fields(map[string]string{
		"method": r.Method,
		"url":    r.URL.String(),
	}).Msg("")
	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"code":   code,
		"status": "error",
		"error":  err.Error(),
	}); err != nil {
		return
	}
}

func JSONResponse(w http.ResponseWriter, code int, data interface{}) {
	if code == http.StatusNoContent {
		w.WriteHeader(code)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"code":   code,
		"status": "success",
		"data":   data,
	}); err != nil {
		return
	}
}
