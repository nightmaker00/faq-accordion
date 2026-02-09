package api

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/nightmaker00/accordion-go/internal/domain"
)

type Handler struct {
	faqService FAQService
}

func NewHandler(faqService FAQService) *Handler {
	return &Handler{faqService: faqService}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/healthz" {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
		return
	}

	const base = "/api/v1/faqs"
	if !strings.HasPrefix(r.URL.Path, base) {
		writeJSON(w, http.StatusNotFound, domain.ErrorResponse{Error: "not found"})
		return
	}

	rest := strings.TrimPrefix(r.URL.Path, base)
	if rest == "" || rest == "/" {
		switch r.Method {
		case http.MethodGet:
			h.handleListFAQs(w, r)
			return
		case http.MethodPost:
			h.handleCreateFAQ(w, r)
			return
		default:
			writeJSON(w, http.StatusMethodNotAllowed, domain.ErrorResponse{Error: "method not allowed"})
			return
		}
	}

	if !strings.HasPrefix(rest, "/") {
		writeJSON(w, http.StatusNotFound, domain.ErrorResponse{Error: "not found"})
		return
	}

	idRaw := strings.TrimPrefix(rest, "/")
	if idRaw == "" || strings.Contains(idRaw, "/") {
		writeJSON(w, http.StatusNotFound, domain.ErrorResponse{Error: "not found"})
		return
	}

	id, err := uuid.Parse(idRaw)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, domain.ErrorResponse{Error: "invalid id"})
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.handleGetFAQ(w, r, id)
		return
	case http.MethodPut:
		h.handleUpdateFAQ(w, r, id)
		return
	case http.MethodDelete:
		h.handleDeleteFAQ(w, r, id)
		return
	default:
		writeJSON(w, http.StatusMethodNotAllowed, domain.ErrorResponse{Error: "method not allowed"})
		return
	}
}

// ListFAQs returns active FAQs ordered by position.
//
// @Summary      List FAQs
// @Description  Get active FAQs ordered by position
// @Tags         faqs
// @Produce      json
// @Success      200  {object}  domain.FAQListResponse
// @Failure      500  {object}  domain.ErrorResponse
// @Router       /faqs [get]
func (h *Handler) handleListFAQs(w http.ResponseWriter, r *http.Request) {
	items, err := h.faqService.ListActive(r.Context())
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, domain.ErrorResponse{Error: "internal error"})
		return
	}

	out := make([]domain.FAQListItemResponse, 0, len(items))
	for _, it := range items {
		out = append(out, domain.FAQListItemResponse{
			ID:       it.ID,
			Title:    it.Title,
			Content:  it.Content,
			Position: it.Position,
		})
	}
	writeJSON(w, http.StatusOK, domain.DataResponse[[]domain.FAQListItemResponse]{Data: out})
}

// GetFAQ returns a FAQ by ID.
//
// @Summary      Get FAQ
// @Description  Get one FAQ by id
// @Tags         faqs
// @Produce      json
// @Param        id   path      string  true  "FAQ ID"
// @Success      200  {object}  domain.FAQResponse
// @Failure      400  {object}  domain.ErrorResponse
// @Failure      404  {object}  domain.ErrorResponse
// @Failure      500  {object}  domain.ErrorResponse
// @Router       /faqs/{id} [get]
func (h *Handler) handleGetFAQ(w http.ResponseWriter, r *http.Request, id uuid.UUID) {
	faq, err := h.faqService.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			writeJSON(w, http.StatusNotFound, domain.ErrorResponse{Error: "not found"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, domain.ErrorResponse{Error: "internal error"})
		return
	}

	writeJSON(w, http.StatusOK, domain.DataResponse[domain.FAQFullResponse]{Data: domain.FAQFullResponse{
		ID:       faq.ID,
		Title:    faq.Title,
		Content:  faq.Content,
		Position: faq.Position,
		IsActive: faq.IsActive,
	}})
}

// CreateFAQ creates a new FAQ.
//
// @Summary      Create FAQ
// @Description  Create new FAQ item
// @Tags         faqs
// @Accept       json
// @Produce      json
// @Param        payload  body      domain.CreateFAQRequest  true  "FAQ payload"
// @Success      201      {object}  domain.FAQResponse
// @Failure      400      {object}  domain.ErrorResponse
// @Failure      500      {object}  domain.ErrorResponse
// @Router       /faqs [post]
func (h *Handler) handleCreateFAQ(w http.ResponseWriter, r *http.Request) {
	var req domain.CreateFAQRequest
	if err := decodeJSON(w, r, &req); err != nil {
		writeJSON(w, http.StatusBadRequest, domain.ErrorResponse{Error: err.Error()})
		return
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	created, err := h.faqService.Create(r.Context(), domain.CreateFAQInput{
		Title:    req.Title,
		Content:  req.Content,
		Position: req.Position,
		IsActive: isActive,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, domain.DataResponse[domain.FAQFullResponse]{Data: domain.FAQFullResponse{
		ID:       created.ID,
		Title:    created.Title,
		Content:  created.Content,
		Position: created.Position,
		IsActive: created.IsActive,
	}})
}

// UpdateFAQ updates a FAQ.
//
// @Summary      Update FAQ
// @Description  Update FAQ by id
// @Tags         faqs
// @Accept       json
// @Produce      json
// @Param        id       path      string                 true  "FAQ ID"
// @Param        payload  body      domain.UpdateFAQRequest true  "FAQ payload"
// @Success      200      {object}  domain.FAQResponse
// @Failure      400      {object}  domain.ErrorResponse
// @Failure      404      {object}  domain.ErrorResponse
// @Failure      500      {object}  domain.ErrorResponse
// @Router       /faqs/{id} [put]
func (h *Handler) handleUpdateFAQ(w http.ResponseWriter, r *http.Request, id uuid.UUID) {
	var req domain.UpdateFAQRequest
	if err := decodeJSON(w, r, &req); err != nil {
		writeJSON(w, http.StatusBadRequest, domain.ErrorResponse{Error: err.Error()})
		return
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	updated, err := h.faqService.Update(r.Context(), id, domain.UpdateFAQInput{
		Title:    req.Title,
		Content:  req.Content,
		Position: req.Position,
		IsActive: isActive,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, domain.DataResponse[domain.FAQFullResponse]{Data: domain.FAQFullResponse{
		ID:       updated.ID,
		Title:    updated.Title,
		Content:  updated.Content,
		Position: updated.Position,
		IsActive: updated.IsActive,
	}})
}

// DeleteFAQ deletes a FAQ.
//
// @Summary      Delete FAQ
// @Description  Delete FAQ by id (idempotent)
// @Tags         faqs
// @Produce      json
// @Param        id   path      string  true  "FAQ ID"
// @Success      200  {object}  domain.MessageResponse
// @Failure      500  {object}  domain.ErrorResponse
// @Router       /faqs/{id} [delete]
func (h *Handler) handleDeleteFAQ(w http.ResponseWriter, r *http.Request, id uuid.UUID) {
	err := h.faqService.Delete(r.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			writeJSON(w, http.StatusNotFound, domain.ErrorResponse{Error: "not found"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, domain.ErrorResponse{Error: "internal error"})
		return
	}
	writeJSON(w, http.StatusOK, domain.MessageResponse{Message: "FAQ deleted successfully"})
}

func writeServiceError(w http.ResponseWriter, err error) {
	if errors.Is(err, domain.ErrNotFound) {
		writeJSON(w, http.StatusNotFound, domain.ErrorResponse{Error: "not found"})
		return
	}

	var ve domain.ValidationError
	if errors.As(err, &ve) {
		writeJSON(w, http.StatusBadRequest, domain.ErrorResponse{Error: ve.Error()})
		return
	}

	writeJSON(w, http.StatusInternalServerError, domain.ErrorResponse{Error: "internal error"})
}

func decodeJSON(w http.ResponseWriter, r *http.Request, dst any) error {
	const maxBytes = 1 << 20 // 1MB
	r.Body = http.MaxBytesReader(w, r.Body, maxBytes)

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(dst); err != nil {
		return err
	}

	// запрещено тащить за собой мусор
	if err := dec.Decode(&struct{}{}); err != io.EOF {
		return errors.New("invalid json")
	}
	return nil
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
