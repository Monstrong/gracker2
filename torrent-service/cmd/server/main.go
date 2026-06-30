package main

import (
	"fmt"
	"log"

	"github.com/monstrong/gracker2/torrent-service/internal/config"
	"github.com/monstrong/gracker2/torrent-service/pkg/db"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("some error loading config")
	}

	// makes pool and does Ping
	pool, err := db.NewPool(&cfg.Postgres)
	if err != nil {
		log.Fatal("some error connecting db")
	}
	defer pool.Close()
	fmt.Println(pool)

	
}
