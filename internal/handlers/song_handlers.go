package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/senyabanana/library-service/internal/entities"
	"github.com/senyabanana/library-service/internal/services"
)

type SongHandler struct {
	service services.SongService
}

func NewSongHandler(service services.SongService) *SongHandler {
	return &SongHandler{service: service}
}

func (c *SongHandler) AddSong(w http.ResponseWriter, r *http.Request) {
	var song entities.Song
	if err := json.NewDecoder(r.Body).Decode(&song); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if err := c.service.AddSong(r.Context(), song); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (c *SongHandler) GetSongs(w http.ResponseWriter, r *http.Request) {
	group := r.URL.Query().Get("group")
	song := r.URL.Query().Get("song")

	page := toInt(r.URL.Query().Get("page"), 1)
	perPage := toInt(r.URL.Query().Get("per_page"), 10)

	filters := entities.SongFilters{
		GroupName: group,
		SongName:  song,
	}
	pagination := entities.Pagination{
		Page:    page,
		PerPage: perPage,
	}

	songs, err := c.service.GetSongs(r.Context(), filters, pagination)
	if err != nil {
		http.Error(w, "Failed to fetch songs", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(songs)
}

func (c *SongHandler) UpdateSong(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/songs/")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var song entities.Song
	if err := json.NewDecoder(r.Body).Decode(&song); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	song.ID = id

	if err := c.service.UpdateSong(r.Context(), song); err != nil {
		http.Error(w, "Failed to update song", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (c *SongHandler) DeleteSong(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/songs/")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if err := c.service.DeleteSong(r.Context(), id); err != nil {
		http.Error(w, "Failed to delete song", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func toInt(value string, defaultValue int) int {
	if i, err := strconv.Atoi(value); err == nil {
		return i
	}
	return defaultValue
}
