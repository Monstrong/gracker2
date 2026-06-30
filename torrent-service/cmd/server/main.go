package main

import (
	"log"

	"github.com/monstrong/gracker2/torrent-service/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("some error loading config")
	}
	
}
