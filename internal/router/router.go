package router

import (
	"net/http"

	_ "github.com/senyabanana/library-service/docs"
	"github.com/senyabanana/library-service/internal/handlers"
	"github.com/senyabanana/library-service/internal/logger"
	"github.com/swaggo/http-swagger"
)

func SetupRoutes(handler *handlers.SongHandler, logg *logger.Logger) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/songs", func(w http.ResponseWriter, r *http.Request) {
		logg.WithField("method", r.Method).Debug("Request received")

		switch r.Method {
		case http.MethodGet:
			handler.GetSongs(w, r)
		case http.MethodPost:
			handler.AddSong(w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/songs/", func(w http.ResponseWriter, r *http.Request) {
		logg.WithField("method", r.Method).Debug("Request received")

		switch r.Method {
		case http.MethodDelete:
			handler.DeleteSong(w, r)
		case http.MethodPut:
			handler.UpdateSong(w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/swagger/", httpSwagger.WrapHandler)

	return mux
}
