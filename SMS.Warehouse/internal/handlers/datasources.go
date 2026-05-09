package handlers

import (
	"database/sql"
	"errors"
	"net/http"

	"sms-warehouse/internal/repository"
)

type DatasourceHandler struct {
	repo   *repository.DatasourceRepo
	dpRepo *repository.DatapointRepo
}

func NewDatasourceHandler(repo *repository.DatasourceRepo, dpRepo *repository.DatapointRepo) *DatasourceHandler {
	return &DatasourceHandler{repo: repo, dpRepo: dpRepo}
}

func (h *DatasourceHandler) List(w http.ResponseWriter, r *http.Request) {
	datasources, err := h.repo.List(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, datasources)
}

func (h *DatasourceHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	ds, err := h.repo.Get(r.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			writeError(w, http.StatusNotFound, "datasource not found")
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	datapoints, err := h.dpRepo.List(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	ds.Datapoints = datapoints

	writeJSON(w, http.StatusOK, ds)
}

func (h *DatasourceHandler) Create(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Type        string `json:"type"`
		Description string `json:"description"`
		Content     string `json:"content"`
	}
	if err := decode(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if body.Type == "" || body.Description == "" {
		writeError(w, http.StatusBadRequest, "type and description are required")
		return
	}

	ds, err := h.repo.Create(r.Context(), body.Type, body.Description, body.Content)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, ds)
}

func (h *DatasourceHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var body struct {
		Type        string `json:"type"`
		Description string `json:"description"`
		Content     string `json:"content"`
	}
	if err := decode(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if body.Type == "" || body.Description == "" {
		writeError(w, http.StatusBadRequest, "type and description are required")
		return
	}

	ds, err := h.repo.Update(r.Context(), id, body.Type, body.Description, body.Content)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			writeError(w, http.StatusNotFound, "datasource not found")
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, ds)
}

func (h *DatasourceHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.repo.Delete(r.Context(), id); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
