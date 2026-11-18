package handlers

import (
	"github.com/monstrong/gracker2/api-gateway/internal/gatewayconfig"
	"github.com/monstrong/gracker2/pkg/logger"
)

type TrackerHandler struct {
    trackerClient TrackerClient
	cfg gatewayconfig.Tracker
	l logger.Logger
}


func NewTrackerHandler(trackerClient TrackerClient, cfg gatewayconfig.Tracker, l logger.Logger) *TrackerHandler {
	return &TrackerHandler{trackerClient: trackerClient, cfg: cfg, l: l}
}
