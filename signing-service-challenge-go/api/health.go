package api

import (
	"net/http"
	"time"
)

type HealthResponse struct {
	Status  string `json:"status"`
	Version string `json:"version"`
}

// Health evaluates the health of the service and writes a standardized response.
func (s *Server) Health(response http.ResponseWriter, request *http.Request) {
	time.Sleep(time.Second * 7)

	health := HealthResponse{
		Status:  "pass",
		Version: "v0",
	}

	WriteAPIResponse(response, http.StatusOK, health)
}
