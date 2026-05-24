package handlers

import (
	"database/sql"
	"errors"
	"net/http"

	"sms-backend/internal/repository"
)

type BoxHandler struct {
	boxes *repository.BoxRepo
}

func NewBoxHandler(boxes *repository.BoxRepo) *BoxHandler {
	return &BoxHandler{boxes: boxes}
}

func (h *BoxHandler) List(w http.ResponseWriter, r *http.Request) {
	pageID := r.PathValue("pageID")
	boxes, err := h.boxes.List(r.Context(), pageID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, boxes)
}

func (h *BoxHandler) Get(w http.ResponseWriter, r *http.Request) {
	pageID := r.PathValue("pageID")
	id := r.PathValue("boxID")
	box, err := h.boxes.Get(r.Context(), pageID, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			writeError(w, http.StatusNotFound, "box not found")
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, box)
}

func (h *BoxHandler) Create(w http.ResponseWriter, r *http.Request) {
	pageID := r.PathValue("pageID")
	var body struct {
		Position     int     `json:"position"`
		DatasourceID *string `json:"datasource_id"`
	}
	if err := decode(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	box, err := h.boxes.Create(r.Context(), pageID, body.Position, body.DatasourceID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, box)
}

func (h *BoxHandler) Update(w http.ResponseWriter, r *http.Request) {
	pageID := r.PathValue("pageID")
	id := r.PathValue("boxID")
	var body struct {
		Position     int     `json:"position"`
		DatasourceID *string `json:"datasource_id"`
	}
	if err := decode(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	box, err := h.boxes.Update(r.Context(), pageID, id, body.Position, body.DatasourceID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			writeError(w, http.StatusNotFound, "box not found")
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, box)
}

func (h *BoxHandler) Delete(w http.ResponseWriter, r *http.Request) {
	pageID := r.PathValue("pageID")
	id := r.PathValue("boxID")
	if err := h.boxes.Delete(r.Context(), pageID, id); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
