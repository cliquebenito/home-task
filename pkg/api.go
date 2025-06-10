package pkg

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"testowoe/internal/models"
	"testowoe/internal/repository"
	"testowoe/internal/service"
)

type Handler struct {
	service *service.BannerService
	log     *slog.Logger
}

type Controller interface {
	RegisterClick(w http.ResponseWriter, r *http.Request)
	CounterView(w http.ResponseWriter, r *http.Request)
	CreateBanner(w http.ResponseWriter, r *http.Request)
}

func NewHandler(serv *service.BannerService, lg *slog.Logger) Controller {
	return &Handler{
		log:     lg,
		service: serv,
	}
}

// RegisterClick registers a click for a banner
func (a *Handler) RegisterClick(w http.ResponseWriter, r *http.Request) {
	const op = "Handler.RegisterClick"

	log := a.log.With(
		slog.String("op", op),
		slog.String("ip", r.RemoteAddr),
	)

	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) < 2 || parts[0] != "counter" {
		log.Warn("invalid path", slog.String("path", r.URL.Path))
		http.Error(w, "invalid path", http.StatusBadRequest)
		return
	}

	idStr := parts[1]
	bannerID, err := strconv.Atoi(idStr)
	if err != nil {
		log.Warn("invalid banner ID", slog.String("banner_id", idStr))
		http.Error(w, "invalid banner ID", http.StatusBadRequest)
		return
	}

	if err := a.service.SaveStatistics(r.Context(), bannerID); err != nil {
		log.Error("invalid banner ID", slog.String("banner_id", idStr))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Info("click registered", slog.Int("banner_id", bannerID))
	w.WriteHeader(http.StatusNoContent)

}

// CounterView returns statistics for a banner
func (a *Handler) CounterView(w http.ResponseWriter, r *http.Request) {
	const op = "Handler.CounterView"
	var req models.RegisterClickRequest

	log := a.log.With(
		slog.String("op", op),
		slog.String("ip", r.RemoteAddr),
	)

	if r.Method != http.MethodPost {
		log.Warn("Method Not Allowed", slog.String("method", r.Method))
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) < 2 || parts[0] != "stats" {
		log.Warn("invalid path", slog.String("path", r.URL.Path))
		http.Error(w, "invalid path", http.StatusBadRequest)
		return
	}

	idStr := parts[1]
	bannerID, err := strconv.Atoi(idStr)
	if err != nil {
		log.Warn("invalid banner ID", slog.String("banner_id", idStr))
		http.Error(w, "invalid banner ID", http.StatusBadRequest)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Warn("invalid JSON body", slog.Any("err", err))
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}
	fromTime, err := parseFlexibleTime(req.From)
	if err != nil {
		log.Warn("invalid 'from' time format", slog.String("from", req.From))
		http.Error(w, "invalid 'from' time format", http.StatusBadRequest)
		return
	}
	toTime, err := parseFlexibleTime(req.To)
	if err != nil {
		log.Warn("invalid 'to' time format", slog.String("to", req.To))
		http.Error(w, "invalid 'to' time format", http.StatusBadRequest)
		return
	}

	stats, err := a.service.LoadStats(r.Context(), models.StatsQueryParams{
		From:     fromTime,
		To:       toTime,
		BannerID: bannerID,
	})
	if err != nil {
		log.Error("failed to load stats", slog.Any("err", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Info("stats loaded", slog.Int("banner_id", bannerID))

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]any{
		"stats": stats,
	}); err != nil {
		log.Error("failed to write response", slog.Any("err", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// CreateBanner creates a new banner
func (a *Handler) CreateBanner(w http.ResponseWriter, r *http.Request) {
	const op = "Handler.CreateBanner"
	var req models.CreateBannerRequest

	log := a.log.With(
		slog.String("op", op),
		slog.String("ip", r.RemoteAddr),
	)

	if r.Method != http.MethodPost {
		log.Warn("Method Not Allowed", slog.String("method", r.Method))
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Warn("invalid JSON body", slog.Any("err", err))
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}

	if err := a.service.CreateBanner(r.Context(), req.Name); err != nil {
		if errors.Is(err, repository.ErrBannerNameExists) {
			log.Warn("banner name already exists", slog.String("name", req.Name))
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(map[string]string{"error": "banner name already exists"})
			return
		}
		log.Error("failed to create banner", slog.String("name", req.Name), slog.Any("err", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Info("banner created", slog.String("name", req.Name))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "created",
	})
}
