package repository

import (
	"context"
	"database/sql"

	"github.com/senyabanana/library-service/internal/entities"
)

type SongRepositoryInterface interface {
	AddSong(ctx context.Context, song entities.Song) error
	GetSongs(ctx context.Context) ([]entities.Song, error)
	GetSongsWithQuery(ctx context.Context, query string, args ...interface{}) ([]entities.Song, error)
	UpdateSong(ctx context.Context, song entities.Song) error
	DeleteSong(ctx context.Context, id int) error
}

type SongRepository struct {
	db *sql.DB
}

func NewSongRepository(db *sql.DB) *SongRepository {
	return &SongRepository{db: db}
}

func (r *SongRepository) AddSong(ctx context.Context, song entities.Song) error {
	query := `INSERT INTO songs (group_name, song_name, release_date, text, link) VALUES ($1, $2, $3, $4, $5)`
	_, err := r.db.ExecContext(ctx, query, song.GroupName, song.SongName, song.ReleaseDate, song.Text, song.Link)
	return err
}

func (r *SongRepository) GetSongs(ctx context.Context) ([]entities.Song, error) {
	query := `SELECT id, group_name, song_name, release_date, text, link FROM songs`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var songs []entities.Song
	for rows.Next() {
		var song entities.Song
		if err := rows.Scan(&song.ID, &song.GroupName, &song.SongName, &song.ReleaseDate, &song.Text, &song.Link); err != nil {
			return nil, err
		}
		songs = append(songs, song)
	}

	return songs, nil
}

func (r *SongRepository) GetSongsWithQuery(ctx context.Context, query string, args ...interface{}) ([]entities.Song, error) {
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var songs []entities.Song
	for rows.Next() {
		var song entities.Song
		if err := rows.Scan(&song.ID, &song.GroupName, &song.SongName, &song.ReleaseDate, &song.Text, &song.Link); err != nil {
			return nil, err
		}
		songs = append(songs, song)
	}

	return songs, nil
}

func (r *SongRepository) UpdateSong(ctx context.Context, song entities.Song) error {
	query := `UPDATE songs SET group_name = $1, song_name = $2, release_date = $3, text = $4, link = $5 WHERE id = $6`
	_, err := r.db.ExecContext(ctx, query, song.GroupName, song.SongName, song.ReleaseDate, song.Text, song.Link, song.ID)
	return err
}

func (r *SongRepository) DeleteSong(ctx context.Context, id int) error {
	query := `DELETE FROM songs WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
