package handlers

import (
	"database/sql"
	"errors"
	"net/http"

	"sms-backend/internal/repository"
)

type PageHandler struct {
	pages *repository.PageRepo
	boxes *repository.BoxRepo
}

func NewPageHandler(pages *repository.PageRepo, boxes *repository.BoxRepo) *PageHandler {
	return &PageHandler{pages: pages, boxes: boxes}
}

func (h *PageHandler) List(w http.ResponseWriter, r *http.Request) {
	pages, err := h.pages.List(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, pages)
}

func (h *PageHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("pageID")
	page, err := h.pages.Get(r.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			writeError(w, http.StatusNotFound, "page not found")
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	boxes, err := h.boxes.List(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	page.Boxes = boxes
	writeJSON(w, http.StatusOK, page)
}

func (h *PageHandler) GetBySlug(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	page, err := h.pages.GetBySlug(r.Context(), slug)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			writeError(w, http.StatusNotFound, "page not found")
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	boxes, err := h.boxes.List(r.Context(), page.ID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	page.Boxes = boxes
	writeJSON(w, http.StatusOK, page)
}

func (h *PageHandler) Create(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Title string `json:"title"`
		Slug  string `json:"slug"`
	}
	if err := decode(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if body.Title == "" || body.Slug == "" {
		writeError(w, http.StatusBadRequest, "title and slug are required")
		return
	}

	page, err := h.pages.Create(r.Context(), body.Title, body.Slug)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, page)
}

func (h *PageHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("pageID")
	var body struct {
		Title string `json:"title"`
		Slug  string `json:"slug"`
	}
	if err := decode(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if body.Title == "" || body.Slug == "" {
		writeError(w, http.StatusBadRequest, "title and slug are required")
		return
	}

	page, err := h.pages.Update(r.Context(), id, body.Title, body.Slug)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			writeError(w, http.StatusNotFound, "page not found")
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, page)
}

func (h *PageHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("pageID")
	if err := h.pages.Delete(r.Context(), id); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
