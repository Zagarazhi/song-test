package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/Zagarazhi/song-test/api"
	"github.com/Zagarazhi/song-test/db"
	"github.com/rs/zerolog/log"
)

func main() {
	db.Init()
	db.Migrate()
	defer db.Close()

	go api.Start()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	log.Info().Msg("Успешное отключение...")
}
