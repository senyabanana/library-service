package services

import (
	"context"
	"errors"
	"strconv"
	"strings"

	//"github.com/senyabanana/library-service/internal/api"
	"github.com/senyabanana/library-service/internal/entities"
	"github.com/senyabanana/library-service/internal/logger"
	"github.com/senyabanana/library-service/internal/repository"
	
	"github.com/sirupsen/logrus"
)

type SongServiceInterface interface {
	AddSong(ctx context.Context, song entities.Song) error
	GetSongs(ctx context.Context, filters entities.SongFilters, pagination entities.Pagination) ([]entities.Song, error)
	GetSongText(ctx context.Context, id int, page int, perPage int) ([]string, error)
	UpdateSong(ctx context.Context, song entities.Song) error
	DeleteSong(ctx context.Context, id int) error
}

type SongService struct {
	repo repository.SongRepositoryInterface
	logg *logger.Logger
}

func NewSongService(repo repository.SongRepositoryInterface, logg *logger.Logger) *SongService {
	return &SongService{
		repo: repo,
		logg: logg,
	}
}

func (s *SongService) AddSong(ctx context.Context, song entities.Song) error {
	s.logg.WithFields(logrus.Fields{
		"group": song.GroupName,
		"song":  song.SongName,
	}).Debug("Adding new song")

	if song.GroupName == "" || song.SongName == "" {
		err := errors.New("group name and song name are required")
		s.logg.WithError(err).Error("Validation failed")
		return err
	}

	//apiClient := api.NewMusicAPIClient("http://external-api")
	//details, err := apiClient.FetchSongDetails(song.GroupName, song.SongName)
	//if err != nil {
	//	return err
	//}

	// Заглушка вместо реального API-запроса
	song.ReleaseDate = "2024-01-01"
	song.Text = "Sample lyrics for " + song.SongName
	song.Link = "https://example.com/" + song.SongName

	err := s.repo.AddSong(ctx, song)
	if err != nil {
		s.logg.WithError(err).Error("Failed to add song to repository")
		return err
	}

	s.logg.WithFields(logrus.Fields{
		"group": song.GroupName,
		"song":  song.SongName,
	}).Info("Song added successfully")
	return nil
}

func (s *SongService) GetSongs(ctx context.Context, filters entities.SongFilters, pagination entities.Pagination) ([]entities.Song, error) {
	s.logg.WithFields(logrus.Fields{
		"filters":    filters,
		"pagination": pagination,
	}).Debug("Fetching songs with filters")

	query := `SELECT id, group_name, song_name, release_date, text, link FROM songs WHERE 1=1`
	args := []interface{}{}
	argIndex := 1

	if filters.GroupName != "" {
		query += ` AND group_name ILIKE $` + strconv.Itoa(argIndex)
		args = append(args, "%"+filters.GroupName+"%")
		argIndex++
	}
	if filters.SongName != "" {
		query += ` AND song_name ILIKE $` + strconv.Itoa(argIndex)
		args = append(args, "%"+filters.SongName+"%")
		argIndex++
	}

	offset := (pagination.Page - 1) * pagination.PerPage
	query += ` LIMIT $` + strconv.Itoa(argIndex) + ` OFFSET $` + strconv.Itoa(argIndex+1)
	args = append(args, pagination.PerPage, offset)

	songs, err := s.repo.GetSongsWithQuery(ctx, query, args...)
	if err != nil {
		s.logg.WithError(err).Error("Failed to fetch songs from repository")
		return nil, err
	}

	s.logg.WithField("count", len(songs)).Info("Songs fetched successfully")
	return songs, nil
}

func (s *SongService) GetSongText(ctx context.Context, id int, page int, perPage int) ([]string, error) {
	s.logg.WithFields(logrus.Fields{
		"song_id": id,
		"page":    page,
		"perPage": perPage,
	}).Debug("Fetching song text")

	songs, err := s.repo.GetSongs(ctx)
	if err != nil {
		s.logg.WithError(err).Error("Failed to fetch songs from repository")
		return nil, err
	}

	var song entities.Song
	for _, s := range songs {
		if s.ID == id {
			song = s
			break
		}
	}

	if song.ID == 0 {
		err := errors.New("song not found")
		s.logg.WithField("song_id", id).Error(err.Error())
		return nil, err
	}

	verses := strings.Split(song.Text, "\n\n")
	start := (page - 1) * perPage
	if start >= len(verses) {
		s.logg.WithFields(logrus.Fields{
			"song_id": id,
			"page":    page,
		}).Debug("No verses found for requested page")
		return nil, nil
	}
	end := start + perPage
	if end > len(verses) {
		end = len(verses)
	}
	s.logg.WithFields(logrus.Fields{
		"song_id": id,
		"page":    page,
	}).Info("Fetched song text successfully")
	return verses[start:end], nil
}

func (s *SongService) UpdateSong(ctx context.Context, song entities.Song) error {
	s.logg.WithFields(logrus.Fields{
		"song_id": song.ID,
		"song":    song.SongName,
	}).Debug("Updating song")

	if song.ID == 0 {
		err := errors.New("song ID is required for update")
		s.logg.WithError(err).Error("Validation failed")
		return err
	}

	err := s.repo.UpdateSong(ctx, song)
	if err != nil {
		s.logg.WithError(err).Error("Failed to update song in repository")
		return err
	}

	s.logg.WithField("song_id", song.ID).Info("Song updated successfully")
	return nil
}

func (s *SongService) DeleteSong(ctx context.Context, id int) error {
	s.logg.WithField("song_id", id).Debug("Deleting song")

	if id == 0 {
		err := errors.New("song ID is required for delete")
		s.logg.WithError(err).Error("Validation failed")
		return err
	}

	err := s.repo.DeleteSong(ctx, id)
	if err != nil {
		s.logg.WithError(err).Error("Failed to delete song from repository")
		return err
	}

	s.logg.WithField("song_id", id).Info("Song deleted successfully")
	return nil
}
