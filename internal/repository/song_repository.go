package repository

import (
	"context"
	"database/sql"

	"github.com/senyabanana/library-service/internal/entities"
	"github.com/senyabanana/library-service/internal/logger"
	
	"github.com/sirupsen/logrus"
)

type SongRepositoryInterface interface {
	AddSong(ctx context.Context, song entities.Song) error
	GetSongs(ctx context.Context) ([]entities.Song, error)
	GetSongsWithQuery(ctx context.Context, query string, args ...interface{}) ([]entities.Song, error)
	UpdateSong(ctx context.Context, song entities.Song) error
	DeleteSong(ctx context.Context, id int) error
}

type SongRepository struct {
	db   *sql.DB
	logg *logger.Logger
}

func NewSongRepository(db *sql.DB, logg *logger.Logger) *SongRepository {
	return &SongRepository{
		db:   db,
		logg: logg,
	}
}

func (r *SongRepository) AddSong(ctx context.Context, song entities.Song) error {
	query := `INSERT INTO songs (group_name, song_name, release_date, text, link) VALUES ($1, $2, $3, $4, $5)`
	r.logg.Debug("Executing query to add song", query)

	_, err := r.db.ExecContext(ctx, query, song.GroupName, song.SongName, song.ReleaseDate, song.Text, song.Link)
	if err != nil {
		r.logg.WithError(err).Error("Failed to execute AddSong query")
		return err
	}

	r.logg.WithField("song", song.SongName).Info("Song added successfully")
	return nil
}

func (r *SongRepository) GetSongs(ctx context.Context) ([]entities.Song, error) {
	query := `SELECT id, group_name, song_name, release_date, text, link FROM songs`
	r.logg.Debug("Executing query to fetch all songs", query)

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		r.logg.WithError(err).Error("Failed to execute GetSongs query")
		return nil, err
	}
	defer rows.Close()

	var songs []entities.Song
	for rows.Next() {
		var song entities.Song
		if err := rows.Scan(&song.ID, &song.GroupName, &song.SongName, &song.ReleaseDate, &song.Text, &song.Link); err != nil {
			r.logg.WithError(err).Error("Failed to scan row in GetSongs")
			return nil, err
		}
		songs = append(songs, song)
	}

	r.logg.WithField("count", len(songs)).Info("Fetched all songs successfully")
	return songs, nil
}

func (r *SongRepository) GetSongsWithQuery(ctx context.Context, query string, args ...interface{}) ([]entities.Song, error) {
	r.logg.WithField("query", query).Debug("Executing query with filters")

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		r.logg.WithError(err).Error("Failed to execute GetSongsWithQuery")
		return nil, err
	}
	defer rows.Close()

	var songs []entities.Song
	for rows.Next() {
		var song entities.Song
		if err := rows.Scan(&song.ID, &song.GroupName, &song.SongName, &song.ReleaseDate, &song.Text, &song.Link); err != nil {
			r.logg.WithError(err).Error("Failed to scan row in GetSongsWithQuery")
			return nil, err
		}
		songs = append(songs, song)
	}

	r.logg.WithField("count", len(songs)).Info("Fetched songs with query successfully")
	return songs, nil
}

func (r *SongRepository) UpdateSong(ctx context.Context, song entities.Song) error {
	query := `UPDATE songs SET group_name = $1, song_name = $2, release_date = $3, text = $4, link = $5 WHERE id = $6`
	r.logg.WithFields(logrus.Fields{
		"query": query,
		"song":  song,
	}).Debug("Executing query to update song")

	_, err := r.db.ExecContext(ctx, query, song.GroupName, song.SongName, song.ReleaseDate, song.Text, song.Link, song.ID)
	if err != nil {
		r.logg.WithError(err).Error("Failed to execute UpdateSong query")
		return err
	}

	r.logg.WithField("song", song.SongName).Info("Song updated successfully")
	return nil
}

func (r *SongRepository) DeleteSong(ctx context.Context, id int) error {
	query := `DELETE FROM songs WHERE id = $1`
	r.logg.WithField("query", query).Debug("Executing query to delete song")

	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		r.logg.WithError(err).Error("Failed to execute DeleteSong query")
		return err
	}

	r.logg.WithField("song_id", id).Info("Song deleted successfully")
	return nil
}
