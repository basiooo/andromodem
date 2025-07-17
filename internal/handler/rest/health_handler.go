package rest

import (
	"net/http"

	"github.com/basiooo/andromodem/internal/common"
)

type HealthHandler struct {
}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) Ping(w http.ResponseWriter, r *http.Request) {
	common.SuccessResponse(w, "pong", nil, http.StatusOK)
}
