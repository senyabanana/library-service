package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/senyabanana/library-service/internal/entities"
	"github.com/senyabanana/library-service/internal/logger"

	"github.com/sirupsen/logrus"
)

type MusicAPIClient struct {
	baseURL string
	logg    *logger.Logger
}

func NewMusicAPIClient(baseURL string, logg *logger.Logger) *MusicAPIClient {
	return &MusicAPIClient{
		baseURL: baseURL,
		logg:    logg,
	}
}

func (c *MusicAPIClient) FetchSongDetails(group, song string) (*entities.Song, error) {
	url := fmt.Sprintf("%s/info?group=%s&song=%s", c.baseURL, group, song)
	c.logg.WithFields(logrus.Fields{
		"url":   url,
		"group": group,
		"song":  song,
	}).Debug("Fetching song details from API")

	resp, err := http.Get(url)
	if err != nil {
		c.logg.WithError(err).Error("Failed to send request to API")
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err := fmt.Errorf("failed to fetch song details: %s", resp.Status)
		c.logg.WithFields(logrus.Fields{
			"url":         url,
			"status_code": resp.StatusCode,
		}).Error(err.Error())
		return nil, err
	}

	var details entities.Details
	if err := json.NewDecoder(resp.Body).Decode(&details); err != nil {
		c.logg.WithError(err).Error("Failed to decode response from API")
		return nil, err
	}

	c.logg.WithFields(logrus.Fields{
		"group":       group,
		"song":        song,
		"releaseDate": details.ReleaseDate,
	}).Info("Fetched song details successfully")

	return &entities.Song{
		GroupName:   group,
		SongName:    song,
		ReleaseDate: details.ReleaseDate,
		Text:        details.Text,
		Link:        details.Link,
	}, nil
}
