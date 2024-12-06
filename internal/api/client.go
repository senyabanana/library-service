package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/senyabanana/library-service/internal/entities"
)

type MusicAPIClient struct {
	baseURL string
}

func NewMusicAPIClient(baseURL string) *MusicAPIClient {
	return &MusicAPIClient{baseURL: baseURL}
}

func (c *MusicAPIClient) FetchSongDetails(group, song string) (*entities.Song, error) {
	url := fmt.Sprintf("%s/info?group=%s&song=%s", c.baseURL, group, song)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch song details: %s", resp.Status)
	}
	
	var details entities.Details
	if err := json.NewDecoder(resp.Body).Decode(&details); err != nil {
		return nil, err
	}

	return &entities.Song{
		GroupName:   group,
		SongName:    song,
		ReleaseDate: details.ReleaseDate,
		Text:        details.Text,
		Link:        details.Link,
	}, nil
}
