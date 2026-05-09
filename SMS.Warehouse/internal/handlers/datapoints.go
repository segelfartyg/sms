package handlers

import (
	"database/sql"
	"errors"
	"net/http"

	"sms-warehouse/internal/repository"
)

type DatapointHandler struct {
	repo *repository.DatapointRepo
}

func NewDatapointHandler(repo *repository.DatapointRepo) *DatapointHandler {
	return &DatapointHandler{repo: repo}
}

func (h *DatapointHandler) List(w http.ResponseWriter, r *http.Request) {
	dsID := r.PathValue("id")
	datapoints, err := h.repo.List(r.Context(), dsID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, datapoints)
}

func (h *DatapointHandler) Get(w http.ResponseWriter, r *http.Request) {
	dsID := r.PathValue("id")
	dpID := r.PathValue("dpID")
	dp, err := h.repo.Get(r.Context(), dsID, dpID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			writeError(w, http.StatusNotFound, "datapoint not found")
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, dp)
}

func (h *DatapointHandler) Create(w http.ResponseWriter, r *http.Request) {
	dsID := r.PathValue("id")
	var body struct {
		Tag         string `json:"tag"`
		Content     string `json:"content"`
		Description string `json:"description"`
	}
	if err := decode(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if body.Tag == "" {
		writeError(w, http.StatusBadRequest, "tag is required")
		return
	}

	dp, err := h.repo.Create(r.Context(), dsID, body.Tag, body.Content, body.Description)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, dp)
}

func (h *DatapointHandler) Update(w http.ResponseWriter, r *http.Request) {
	dsID := r.PathValue("id")
	dpID := r.PathValue("dpID")
	var body struct {
		Tag         string `json:"tag"`
		Content     string `json:"content"`
		Description string `json:"description"`
	}
	if err := decode(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if body.Tag == "" {
		writeError(w, http.StatusBadRequest, "tag is required")
		return
	}

	dp, err := h.repo.Update(r.Context(), dsID, dpID, body.Tag, body.Content, body.Description)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			writeError(w, http.StatusNotFound, "datapoint not found")
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, dp)
}

func (h *DatapointHandler) Delete(w http.ResponseWriter, r *http.Request) {
	dsID := r.PathValue("id")
	dpID := r.PathValue("dpID")
	if err := h.repo.Delete(r.Context(), dsID, dpID); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
