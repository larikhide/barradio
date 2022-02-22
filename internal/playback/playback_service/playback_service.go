package playback_service

import (
	"context"
	"errors"

	"github.com/zmb3/spotify/v2"
	sauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2/clientcredentials"
)

var (
	ErrNotFound             = errors.New("search result is empty")
	ErrUnknownCategory      = errors.New("unknown category")
	CategoriesSearchStrings = map[string]string{
		"relaxed":  "relax",
		"lyrical":  "lyric",
		"cheerful": "cheerful",
	}
)

// PlaybackService provides methods for work with music stream service
type PlaybackService struct {
	clientID     string
	clientSecret string
}

func NewPlaybackService(clientID, clientSecret string) (*PlaybackService, error) {
	return &PlaybackService{
		clientID:     clientID,
		clientSecret: clientSecret,
	}, nil
}

func (s *PlaybackService) getClient() (*spotify.Client, error) {
	ctx := context.Background()
	config := &clientcredentials.Config{
		ClientID:     s.clientID,
		ClientSecret: s.clientSecret,
		TokenURL:     sauth.TokenURL,
	}
	token, err := config.Token(ctx)
	if err != nil {
		return nil, err
	}

	httpClient := sauth.New().Client(ctx, token)
	client := spotify.New(httpClient)
	return client, nil
}

// GetTracklistByCategory fetches a playlist related to defined category.
//
// It is very simple now: search playlist by hardcoded key-word,
// stored as CategoriesSearchStrings map
func (s *PlaybackService) GetTracklistByCategory(category string) (*TrackList, error) {
	searchText, ok := CategoriesSearchStrings[category]
	if !ok {
		return nil, ErrUnknownCategory
	}

	client, err := s.getClient()
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	searchResults, err := client.Search(ctx, searchText, spotify.SearchTypePlaylist)
	if err != nil {
		return nil, err
	}

	if len(searchResults.Playlists.Playlists) == 0 {
		return nil, ErrNotFound
	}

	playlist := searchResults.Playlists.Playlists[0]
	plTracks, err := client.GetPlaylistTracks(ctx, playlist.ID)
	if err != nil {
		return nil, err
	}

	tracks := make([]*Track, 0, len(plTracks.Tracks))
	totalDuration := 0.0
	for _, tr := range plTracks.Tracks {
		// MVP: fetch first items only
		artist := ""
		if len(tr.Track.Artists) > 0 {
			artist = tr.Track.Artists[0].Name
		}
		image := ""
		if len(tr.Track.Album.Images) > 0 {
			image = tr.Track.Album.Images[0].URL
		}
		url := ""
		if len(tr.Track.ExternalURLs) > 0 {
			for _, val := range tr.Track.ExternalURLs {
				url = val
				break
			}
		}
		tracks = append(tracks, &Track{
			ID:          string(tr.Track.ID),
			Name:        tr.Track.Name,
			Album:       tr.Track.Album.Name,
			Artist:      artist,
			ImageURL:    image,
			ExternalURL: url,
			Duration:    tr.Track.TimeDuration().Seconds(),
		})
		totalDuration += tr.Track.TimeDuration().Seconds()
	}

	url := ""
	if len(playlist.ExternalURLs) > 0 {
		for _, val := range playlist.ExternalURLs {
			url = val
			break
		}
	}
	result := &TrackList{
		ID:            string(playlist.ID),
		Name:          playlist.Name,
		ExternalURL:   url,
		Tracks:        tracks,
		TotalTracks:   len(tracks),
		TotalDuration: totalDuration,
	}
	return result, nil
}

// GetTracksByCategory fetches a few tracks related to defined category.
//
// It is very simple now: search playlist by hardcoded key-word and fetch all tracks from it.
// key-words stored as CategoriesSearchStrings map
func (s *PlaybackService) GetTracksByCategory(category string) ([]*Track, error) {
	playlist, err := s.GetTracklistByCategory(category)
	if err != nil && errors.Is(err, ErrNotFound) {
		return []*Track{}, nil
	}
	if err != nil {
		return nil, err
	}

	return playlist.Tracks, nil
}
