package db

import (
	"fmt"
	"os"

	"github.com/Zagarazhi/song-test/models"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Сохраненное подключение к БД
var (
	connection *gorm.DB
)

// Инициализация подключения
func Init() {
	log.Info().Msg("Получение данных из .env файла")
	err := godotenv.Load()
	if err != nil {
		log.Fatal().Msg("Не удалось загрузить .env файл")
	}

	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")

	dsn := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable", dbUser, dbPassword, dbName, dbHost, dbPort)
	conn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})

	if err != nil {
		log.Fatal().Msg("Не удалось подключиться к базе данных")
	}

	connection = conn
}

// Получение подключения
func GetConnection() *gorm.DB {
	return connection
}

// Закрытие подключения
func Close() {
	db, err := connection.DB()
	if err != nil {
		log.Error().Err(err).Send()
	}

	if err := db.Close(); err != nil {
		log.Error().Err(err).Send()
	}

	log.Info().Msg("Подключение к БД успешно закрыто")
}

// Автомиграция
func Migrate() {
	if err := connection.AutoMigrate(
		&models.SongGorm{},
	); err != nil {
		log.Error().Err(err).Send()
	}
}
