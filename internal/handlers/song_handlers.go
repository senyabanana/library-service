package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/senyabanana/library-service/internal/entities"
	"github.com/senyabanana/library-service/internal/logger"
	"github.com/senyabanana/library-service/internal/services"

	"github.com/sirupsen/logrus"
)

type SongHandler struct {
	service services.SongServiceInterface
	logg    *logger.Logger
}

func NewSongHandler(service services.SongServiceInterface, logg *logger.Logger) *SongHandler {
	return &SongHandler{
		service: service,
		logg:    logg,
	}
}

func (h *SongHandler) AddSong(w http.ResponseWriter, r *http.Request) {
	h.logg.WithField("method", r.Method).Debug("Handling AddSong request")

	var song entities.Song
	if err := json.NewDecoder(r.Body).Decode(&song); err != nil {
		h.logg.WithError(err).Error("Invalid request payload")
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if err := h.service.AddSong(r.Context(), song); err != nil {
		h.logg.WithError(err).Error("Failed to add song")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.logg.WithFields(logrus.Fields{
		"group": song.GroupName,
		"song":  song.SongName,
	}).Info("Song added successfully")

	w.WriteHeader(http.StatusCreated)
}

func (h *SongHandler) GetSongs(w http.ResponseWriter, r *http.Request) {
	h.logg.WithField("method", r.Method).Debug("Handling GetSongs request")

	group := r.URL.Query().Get("group")
	song := r.URL.Query().Get("song")

	page := toInt(r.URL.Query().Get("page"), 1)
	perPage := toInt(r.URL.Query().Get("per_page"), 10)

	h.logg.WithFields(logrus.Fields{
		"group":    group,
		"song":     song,
		"page":     page,
		"per_page": perPage,
	}).Debug("Query parameters for GetSongs")

	filters := entities.SongFilters{
		GroupName: group,
		SongName:  song,
	}
	pagination := entities.Pagination{
		Page:    page,
		PerPage: perPage,
	}

	songs, err := h.service.GetSongs(r.Context(), filters, pagination)
	if err != nil {
		h.logg.WithError(err).Error("Failed to fetch songs")
		http.Error(w, "Failed to fetch songs", http.StatusInternalServerError)
		return
	}

	h.logg.WithField("count", len(songs)).Info("Fetched songs successfully")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(songs)
}

func (h *SongHandler) UpdateSong(w http.ResponseWriter, r *http.Request) {
	h.logg.WithField("method", r.Method).Debug("Handling UpdateSong request")

	idStr := strings.TrimPrefix(r.URL.Path, "/songs/")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		h.logg.WithField("id", idStr).Error("Invalid ID")
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var song entities.Song
	if err := json.NewDecoder(r.Body).Decode(&song); err != nil {
		h.logg.WithError(err).Error("Invalid request payload")
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	song.ID = id

	if err := h.service.UpdateSong(r.Context(), song); err != nil {
		h.logg.WithError(err).WithField("id", id).Error("Failed to update song")
		http.Error(w, "Failed to update song", http.StatusInternalServerError)
		return
	}

	h.logg.WithField("id", id).Info("Song updated successfully")
	w.WriteHeader(http.StatusOK)
}

func (h *SongHandler) DeleteSong(w http.ResponseWriter, r *http.Request) {
	h.logg.WithField("method", r.Method).Debug("Handling DeleteSong request")

	idStr := strings.TrimPrefix(r.URL.Path, "/songs/")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		h.logg.WithField("id", idStr).Error("Invalid ID")
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if err := h.service.DeleteSong(r.Context(), id); err != nil {
		h.logg.WithError(err).WithField("id", id).Error("Failed to delete song")
		http.Error(w, "Failed to delete song", http.StatusInternalServerError)
		return
	}

	h.logg.WithField("id", id).Info("Song deleted successfully")
	w.WriteHeader(http.StatusNoContent)
}

func toInt(value string, defaultValue int) int {
	if i, err := strconv.Atoi(value); err == nil {
		return i
	}
	return defaultValue
}
