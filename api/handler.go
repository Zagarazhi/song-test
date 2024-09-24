package api

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/Zagarazhi/song-test/db"
	"github.com/Zagarazhi/song-test/models"
	"github.com/rs/zerolog/log"
	"github.com/valyala/fasthttp"
)

// Получение песен с фильтрацией и пагинацией
func fetchSongs(ctx *fasthttp.RequestCtx) {
	id := ctx.QueryArgs().GetUintOrZero("id")
	group := string(ctx.QueryArgs().Peek("group"))
	song := string(ctx.QueryArgs().Peek("song"))
	text := string(ctx.QueryArgs().Peek("text"))
	link := string(ctx.QueryArgs().Peek("link"))
	startTimeStr := string(ctx.QueryArgs().Peek("startTime"))
	endTimeStr := string(ctx.QueryArgs().Peek("endTime"))
	offset := ctx.QueryArgs().GetUintOrZero("offset")
	limit := ctx.QueryArgs().GetUintOrZero("limit")

	if limit == 0 {
		limit = 10
	}

	var startTime, endTime time.Time
	var err error
	layout := "02.01.2006"
	if len(startTimeStr) > 0 {
		startTime, err = time.Parse(layout, startTimeStr)
		if err != nil {
			OutputJsonMessageResult(ctx, 400, "Не могу конвертировать startTime")
			return
		}
	}
	if len(endTimeStr) > 0 {
		endTime, err = time.Parse(layout, endTimeStr)
		if err != nil {
			OutputJsonMessageResult(ctx, 400, "Не могу конвертировать endTime")
			return
		}
	}

	songs, err := db.GetFullInfos(uint64(id), group, song, text, link, startTime, endTime, uint64(offset), uint64(limit))
	if err != nil {
		OutputJson(ctx, 500, "Что-то пошло не так")
		return
	}
	OutputJsonNoIndent(ctx, 200, ConvertGormSongsToSongs(songs))
}

// Метод получения текста песен с пагинацией
func fetchText(ctx *fasthttp.RequestCtx) {
	id := ctx.QueryArgs().GetUintOrZero("id")
	offset := ctx.QueryArgs().GetUintOrZero("offset")
	limit := ctx.QueryArgs().GetUintOrZero("limit")

	if id == 0 {
		OutputJson(ctx, 400, "id не может быть нулем")
	}

	if limit == 0 {
		limit = 10
	}

	text, err := db.GetSongText(uint64(id), uint64(offset), uint64(limit))
	if err != nil {
		OutputJson(ctx, 500, "Что-то пошло не так")
		return
	}
	OutputJsonNoIndent(ctx, 200, text)
}

// Метод получения деталей песни
func fetchDetails(ctx *fasthttp.RequestCtx) {
	group := string(ctx.QueryArgs().Peek("group"))
	song := string(ctx.QueryArgs().Peek("song"))

	details, err := db.GetSongDetails(group, song)
	if err != nil {
		OutputJson(ctx, 500, "Что-то пошло не так")
		return
	}
	if details.ID == 0 {
		OutputJson(ctx, 404, "Запись не найдена")
		return
	}
	OutputJsonNoIndent(ctx, 200, ConvertSongGormToDetails(details))
}

// Вставка песен
func insertSongs(ctx *fasthttp.RequestCtx) {
	reqBody := ctx.Request.Body()
	var reqParams []models.AddSong

	if err := json.Unmarshal(reqBody, &reqParams); err != nil {
		err = errors.New("unmarshal: " + err.Error())
		log.Error().Err(err).Send()
		OutputJsonMessageResult(ctx, 400, err.Error())
		return
	}

	res, err := db.CreateSong(ConvertAddSongsToGormSongs(reqParams))

	if err != nil {
		OutputJsonMessageResult(ctx, 500, "Что-то пошло не так")
		return
	}

	OutputJsonNoIndent(ctx, 200, ConvertGormSongsToSongs(res))
}

// Обновление песни
func updateSong(ctx *fasthttp.RequestCtx) {
	reqBody := ctx.Request.Body()
	var reqParams models.Song

	if err := json.Unmarshal(reqBody, &reqParams); err != nil {
		err = errors.New("unmarshal: " + err.Error())
		log.Error().Err(err).Send()
		OutputJsonMessageResult(ctx, 400, err.Error())
		return
	}

	res, err := db.UpdateSong(ConvertSongToSongGorm(reqParams))

	if err != nil {
		OutputJsonMessageResult(ctx, 500, "Что-то пошло не так")
		return
	}

	OutputJsonNoIndent(ctx, 200, ConvertSongGormToSong(res))
}

// Удаление песни
func deleteSong(ctx *fasthttp.RequestCtx) {
	id := ctx.QueryArgs().GetUintOrZero("id")

	if id == 0 {
		OutputJson(ctx, 400, "id не может быть нулем")
	}

	res, err := db.DeleteSong(uint64(id))
	if err != nil {
		OutputJsonMessageResult(ctx, 500, "Что-то пошло не так")
		return
	}
	OutputJsonNoIndent(ctx, 200, res)
}

func options(ctx *fasthttp.RequestCtx) {
	OutputCORSOptions(ctx)
}
