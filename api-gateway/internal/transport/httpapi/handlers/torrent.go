package handlers

import (
	"github.com/monstrong/gracker2/api-gateway/internal/gatewayconfig"
	"github.com/monstrong/gracker2/pkg/logger"
)

type TorrentHandler struct {
    torrentClient TorrentClient
	cfg gatewayconfig.Torrent
	l logger.Logger
}


func NewTorrentHandler(torrentClient TorrentClient, cfg gatewayconfig.Torrent, l logger.Logger) *TorrentHandler {
	return &TorrentHandler{
		torrentClient: torrentClient, 
		cfg: cfg, 
		l: l,
	}
}
