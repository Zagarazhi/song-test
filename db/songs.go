package db

import (
	"context"
	"strings"
	"time"

	"github.com/Zagarazhi/song-test/models"
	"github.com/rs/zerolog/log"
)

// GetFullInfos возвращает список песен с фильтрацией по полям и пагинацией
func GetFullInfos(id uint64, group, song, text, link string, startTime, endTime time.Time, offset, limit uint64) ([]models.SongGorm, error) {
	// Подключение к базе данных
	connection := GetConnection()

	// Создаем контекст с таймаутом для запроса
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Начинаем запрос с модели SongGorm с использованием контекста
	query := connection.Model(&models.SongGorm{}).WithContext(ctx)

	// Добавляем фильтр по id, если передан
	if id > 0 {
		query = query.Where("id = ?", id)
	}

	// Добавляем фильтр по group, если передан
	if len(group) > 0 {
		query = query.Where("group_name ILIKE ?", "%"+group+"%")
	}

	// Добавляем фильтр по song, если передан
	if len(song) > 0 {
		query = query.Where("song ILIKE ?", "%"+song+"%")
	}

	// Добавляем фильтр по text, если передан
	if len(text) > 0 {
		query = query.Where("text ILIKE ?", "%"+text+"%")
	}

	// Добавляем фильтр по link, если передан
	if len(link) > 0 {
		query = query.Where("link ILIKE ?", "%"+link+"%")
	}

	// Фильтр по дате создания
	if !startTime.IsZero() && !endTime.IsZero() {
		// Оба параметра заданы: фильтр по диапазону
		query = query.Where("release_date BETWEEN ? AND ?", startTime, endTime)
	} else if !startTime.IsZero() {
		// Только startTime задан: фильтр с даты начала до настоящего времени
		query = query.Where("release_date >= ?", startTime)
	} else if !endTime.IsZero() {
		// Только endTime задан: фильтр от начала до указанной даты
		query = query.Where("release_date <= ?", endTime)
	}

	// Добавляем пагинацию
	query = query.Offset(int(offset)).Limit(int(limit))

	// Выполняем запрос
	var songs []models.SongGorm
	if err := query.Find(&songs).Error; err != nil {
		log.Error().Err(err).Send()
		return nil, err
	}

	return songs, nil
}

// Получение текста песни по куплетам
func GetSongText(id, offset, limit uint64) ([]string, error) {
	// Подключение к базе данных
	connection := GetConnection()

	// Создаем контекст с таймаутом для запроса
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Инициализируем переменную для хранения песни
	var song models.SongGorm
	verses := make([]string, 0)

	// Выполняем запрос на получение текста песни по id
	row := connection.WithContext(ctx).Model(&models.SongGorm{}).Where("id = ?", id).First(&song).RowsAffected
	if row < 1 {
		return verses, nil
	}

	// Разбиваем текст песни на массив строк по разделителю (перенос строки)
	verses = strings.Split(song.Text, "\n\n")

	// Применяем пагинацию
	start := int(offset)
	end := start + int(limit)

	// Если конец выходит за границы массива, корректируем его
	if end > len(verses) {
		end = len(verses)
	}

	// Если offset больше длины массива, возвращаем пустой результат
	if start > len(verses) {
		return []string{}, nil
	}

	// Возвращаем куплеты с учетом пагинации
	return verses[start:end], nil
}

// Получение информации по деталям песни из группы и названия
func GetSongDetails(group, song string) (models.SongGorm, error) {
	// Подключение к базе данных
	connection := GetConnection()

	// Создаем контекст с таймаутом для запроса
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Инициализируем переменную для хранения песни
	var res models.SongGorm

	// Запрос с модели SongGorm с использованием контекста
	connection.Model(&models.SongGorm{}).WithContext(ctx).Where("group_name = ? AND song = ?", group, song).First(&res)
	return res, nil
}

// Создание песни в базе
func CreateSong(songs []models.SongGorm) ([]models.SongGorm, error) {
	// Подключение к базе данных
	connection := GetConnection()

	// Создаем контекст с таймаутом для запроса
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Создание песен в батче
	err := connection.Model(&models.SongGorm{}).WithContext(ctx).CreateInBatches(&songs, 500).Error
	if err != nil {
		log.Error().Err(err).Send()
		return nil, err
	}
	return songs, nil
}

// Обновление песни по ID
func UpdateSong(song models.SongGorm) (models.SongGorm, error) {
	// Подключение к базе данных
	connection := GetConnection()

	// Создаем контекст с таймаутом для запроса
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := connection.Model(&models.SongGorm{}).WithContext(ctx).Where("id = ?", song.ID).Updates(&song).Error
	if err != nil {
		log.Error().Err(err).Send()
		return song, err
	}
	return song, nil
}

// Удаление по ID
func DeleteSong(id uint64) (uint64, error) {
	// Подключение к базе данных
	connection := GetConnection()

	// Создаем контекст с таймаутом для запроса
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := connection.Model(&models.SongGorm{}).WithContext(ctx).Delete(&models.SongGorm{}, id).Error
	if err != nil {
		log.Error().Err(err).Send()
		return 0, err
	}
	return id, nil
}
