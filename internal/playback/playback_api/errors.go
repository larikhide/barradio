package playback_api

// APIBaseError contains api error message to return to client
type APIBaseError struct {
	Message string `json:"message"`
}
