package api

import (
	"os"

	"github.com/fasthttp/router"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
	"github.com/valyala/fasthttp"
)

// Функция запуска сервера
func Start() {
	log.Info().Msg("Получение данных из .env файла")
	err := godotenv.Load()
	if err != nil {
		log.Fatal().Msg("Не удалось загрузить .env файл")
	}
	host := os.Getenv("HOST")

	log.Info().Msgf("Запуск апи на хосту %s", host)

	r := router.New()

	r.GET("/songs", fetchSongs)
	r.POST("/songs", insertSongs)
	r.PUT("/songs", updateSong)
	r.DELETE("/songs", deleteSong)
	r.OPTIONS("/songs", options)

	r.GET("/info", fetchDetails)
	r.OPTIONS("/info", options)

	r.GET("/text", fetchText)
	r.OPTIONS("/text", options)

	server := &fasthttp.Server{
		Handler:            r.Handler,
		MaxRequestBodySize: 1024 * 1024,
	}

	log.Fatal().Err(server.ListenAndServe(host)).Send()
}
