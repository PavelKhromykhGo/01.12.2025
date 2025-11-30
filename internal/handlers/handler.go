package handlers

import (
	"LinkChecker/internal/models"
	"LinkChecker/internal/service"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type PDFGenerator interface {
	Generate(links []models.LinksGroup) ([]byte, error)
}

type Handler struct {
	svc    service.Service
	pdfGen PDFGenerator
}

func NewHandler(s service.Service, p PDFGenerator) *Handler {
	return &Handler{svc: s, pdfGen: p}
}

func (h *Handler) Route() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)

	r.Post("/links/check", h.CheckLinks)
	r.Post("/links/report", h.Report)

	return r
}

type checkLinksRequest struct {
	Links []string `json:"links"`
}

type checkLinksResponse struct {
	Links    map[string]string `json:"links"`
	LinksNum int               `json:"links_num"`
}

func (h *Handler) CheckLinks(w http.ResponseWriter, r *http.Request) {
	var req checkLinksRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	if len(req.Links) == 0 {
		http.Error(w, "no links", http.StatusBadRequest)
		return
	}

	group, err := h.svc.CheckLinks(r.Context(), req.Links)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	resp := checkLinksResponse{
		Links:    make(map[string]string, len(group.Links)),
		LinksNum: group.ID,
	}

	for _, link := range group.Links {
		resp.Links[link.URL] = link.Status
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

type ReportRequest struct {
	LinksList []int `json:"links_list"`
}

func (h *Handler) Report(w http.ResponseWriter, r *http.Request) {
	var req ReportRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	if len(req.LinksList) == 0 {
		http.Error(w, "links_list empty", http.StatusBadRequest)
		return
	}

	groups, err := h.svc.GetGroups(r.Context(), req.LinksList)
	if err != nil {
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}
	if len(groups) == 0 {
		http.Error(w, "no groups", http.StatusNotFound)
	}

	pdf, err := h.pdfGen.Generate(groups)
	if err != nil {
		http.Error(w, "pdf generate error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "attachment; filename=\"report.pdf\"")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(pdf)
}
