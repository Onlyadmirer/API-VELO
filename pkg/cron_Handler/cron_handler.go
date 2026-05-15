package cronhandler

import (
	"VELO-backend/pkg/utils"
	"net/http"
)

func CronSendResponse(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	utils.ResponseSuccess(w, http.StatusOK, "status: ok", nil)
}
