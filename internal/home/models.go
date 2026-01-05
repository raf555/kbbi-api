package home

import "github.com/raf555/kbbi-api/internal/dictionary"

type HomeResponse struct {
	Message       string           `json:"message"`
	Stats         dictionary.Stats `json:"stats"`
	Documentation string           `json:"documentation"`
	Issues        string           `json:"issues"`
}

type HealthResponse struct {
	Message string `json:"message"`
}
