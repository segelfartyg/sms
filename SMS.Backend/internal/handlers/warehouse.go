package handlers

import (
	"net/http"

	"sms-backend/internal/warehouse"
)

type WarehouseHandler struct {
	client *warehouse.Client
}

func NewWarehouseHandler(client *warehouse.Client) *WarehouseHandler {
	return &WarehouseHandler{client: client}
}

func (h *WarehouseHandler) ListDatasources(w http.ResponseWriter, r *http.Request) {
	datasources, err := h.client.ListDatasources(r.Context())
	if err != nil {
		writeError(w, http.StatusBadGateway, "warehouse unavailable: "+err.Error())
		return
	}
	writeJSON(w, http.StatusOK, datasources)
}
