package services

import (
	"context"
	"errors"
	"strconv"
	"strings"

	//"github.com/senyabanana/library-service/internal/api"
	"github.com/senyabanana/library-service/internal/entities"
	"github.com/senyabanana/library-service/internal/repository"
)

type SongServiceInterface interface {
	AddSong(ctx context.Context, song entities.Song) error
	GetSongs(ctx context.Context) ([]entities.Song, error)
	GetSongText(ctx context.Context, id int, page int, perPage int) ([]string, error)
	UpdateSong(ctx context.Context, song entities.Song) error
	DeleteSong(ctx context.Context, id int) error
}

type SongService struct {
	repo repository.SongRepositoryInterface
}

func NewSongService(repo repository.SongRepositoryInterface) *SongService {
	return &SongService{repo: repo}
}

func (s *SongService) AddSong(ctx context.Context, song entities.Song) error {
	if song.GroupName == "" || song.SongName == "" {
		return errors.New("group name and song name are required")
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

	return s.repo.AddSong(ctx, song)
}

func (s *SongService) GetSongs(ctx context.Context, filters entities.SongFilters, pagination entities.Pagination) ([]entities.Song, error) {
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

	return s.repo.GetSongsWithQuery(ctx, query, args...)
}

func (s *SongService) GetSongText(ctx context.Context, id int, page int, perPage int) ([]string, error) {
	songs, err := s.repo.GetSongs(ctx)
	if err != nil {
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
		return nil, errors.New("song not found")
	}

	verses := strings.Split(song.Text, "\n\n")
	start := (page - 1) * perPage
	if start >= len(verses) {
		return nil, nil
	}
	end := start + perPage
	if end > len(verses) {
		end = len(verses)
	}
	return verses[start:end], nil
}

func (s *SongService) UpdateSong(ctx context.Context, song entities.Song) error {
	if song.ID == 0 {
		return errors.New("song ID is required for update")
	}
	return s.repo.UpdateSong(ctx, song)
}

func (s *SongService) DeleteSong(ctx context.Context, id int) error {
	if id == 0 {
		return errors.New("song ID is required for delete")
	}
	return s.repo.DeleteSong(ctx, id)
}
