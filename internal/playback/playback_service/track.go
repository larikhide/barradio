package playback_service

type Track struct {
	ID          string
	Name        string
	Album       string
	Artist      string
	ImageURL    string
	ExternalURL string
	Duration    float64
}

type TrackList struct {
	ID            string
	Name          string
	ExternalURL   string
	Tracks        []*Track
	TotalTracks   int
	TotalDuration float64
}
