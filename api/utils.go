package api

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/Zagarazhi/song-test/models"
	"github.com/rs/zerolog/log"
	"github.com/valyala/fasthttp"
)

type out struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func OutputJsonMessageResult(ctx *fasthttp.RequestCtx, code int, r string) {
	ctx.Response.Header.Set("Content-Type", "application/json")
	ctx.Response.Header.Set("Access-Control-Allow-Origin", "*")
	ctx.Response.Header.Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	ctx.Response.Header.Set("Access-Control-Allow-Headers", "Authorization")
	ctx.Response.Header.SetStatusCode(code)
	out := out{code, r}
	jsonResult, _ := json.Marshal(out)
	if _, err := fmt.Fprint(ctx, string(jsonResult)); err != nil {
		log.Error().Err(err).Send()
	}
	ctx.Response.Header.Set("Connection", "close")
}

func OutputCORSOptions(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.Set("Content-Type", "text/html")
	ctx.Response.Header.Set("Access-Control-Allow-Origin", "*")
	ctx.Response.Header.Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	ctx.Response.Header.Set("Access-Control-Allow-Headers", "*")
	ctx.Response.Header.SetStatusCode(200)
	ctx.Response.Header.Set("Connection", "close")
}

func OutputJson(ctx *fasthttp.RequestCtx, code int, result interface{}) {
	jsonResult, err := json.MarshalIndent(result, "", "    ")
	if err != nil {
		log.Error().Err(err).Send()
		OutputJsonMessageResult(ctx, 500, "errors.common.internalError")
		return
	}
	ctx.Response.Header.Set("Content-Type", "application/json")
	ctx.Response.Header.Set("Access-Control-Allow-Origin", "*")
	ctx.Response.Header.Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	ctx.Response.Header.Set("Access-Control-Allow-Headers", "Authorization")
	ctx.Response.SetStatusCode(code)
	if _, err = fmt.Fprint(ctx, string(jsonResult)); err != nil {
		log.Error().Err(err).Send()
	}
	ctx.Response.Header.Set("Connection", "close")
}

func OutputJsonNoIndent(ctx *fasthttp.RequestCtx, code int, result interface{}) {
	jsonResult, err := json.Marshal(result)
	if err != nil {
		log.Error().Err(err).Send()
		OutputJsonMessageResult(ctx, 500, "errors.common.internalError")
		return
	}
	ctx.Response.Header.Set("Content-Type", "application/json")
	ctx.Response.Header.Set("Access-Control-Allow-Origin", "*")
	ctx.Response.Header.Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	ctx.Response.Header.Set("Access-Control-Allow-Headers", "Authorization")
	ctx.Response.SetStatusCode(code)
	if _, err = fmt.Fprint(ctx, string(jsonResult)); err != nil {
		log.Error().Err(err).Send()
	}
	ctx.Response.Header.Set("Connection", "close")
}

// Конвертируем хранимый формат в формат вывода
func ConvertGormSongsToSongs(songs []models.SongGorm) []models.Song {
	layout := "02.01.2006"
	res := make([]models.Song, len(songs))
	for i, song := range songs {
		temp := models.Song{
			ID:          song.ID,
			Group:       song.GroupName,
			Song:        song.Song,
			Text:        song.Text,
			Link:        song.Link,
			ReleaseDate: song.ReleaseDate.Format(layout),
		}
		res[i] = temp
	}
	return res
}

// Конвертируем песни для вставки в формат модели БД
func ConvertAddSongsToGormSongs(songs []models.AddSong) []models.SongGorm {
	res := make([]models.SongGorm, len(songs))
	for i, song := range songs {
		temp := models.SongGorm{
			GroupName: song.Group,
			Song:      song.Song,
		}
		res[i] = temp
	}
	return res
}

// Конвертируем песню из БД в детальную информацию
func ConvertSongGormToDetails(song models.SongGorm) models.SongDetails {
	layout := "02.01.2006"
	return models.SongDetails{
		ReleaseDate: song.ReleaseDate.Format(layout),
		Text:        song.Text,
		Link:        song.Link,
	}
}

// Конвертируем песню в формат БД
func ConvertSongGormToSong(song models.SongGorm) models.Song {
	layout := "02.01.2006"
	return models.Song{
		ID:          song.ID,
		Group:       song.GroupName,
		Song:        song.Song,
		Text:        song.Text,
		Link:        song.Link,
		ReleaseDate: song.ReleaseDate.Format(layout),
	}
}

// Конвертируем песню из формата БД
func ConvertSongToSongGorm(song models.Song) models.SongGorm {
	layout := "02.01.2006"
	t := time.Time{}
	if len(song.ReleaseDate) > 0 {
		t, _ = time.Parse(layout, song.ReleaseDate)
	}
	return models.SongGorm{
		ID:          song.ID,
		GroupName:   song.Group,
		Song:        song.Song,
		Text:        song.Text,
		Link:        song.Link,
		ReleaseDate: t,
	}
}
