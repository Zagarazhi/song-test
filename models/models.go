package models

import "time"

// Модель для хранения песни в базе
type SongGorm struct {
	ID          uint64 `gorm:"primaryKey;autoIncrement"`
	GroupName   string `gorm:"type:varchar(255);not null"`
	Song        string `gorm:"type:varchar(255);not null"`
	ReleaseDate time.Time
	Text        string `gorm:"type:text"`
	Link        string `gorm:"type:text"`
}

func (SongGorm) TableName() string {
	return "songs"
}

// Модель для вывода ее пользователю
type Song struct {
	ID          uint64 `json:"id"`
	Group       string `json:"group"`
	Song        string `json:"song"`
	ReleaseDate string `json:"releaseDate"`
	Text        string `json:"text"`
	Link        string `json:"link"`
}

// Детали песни
type SongDetails struct {
	ReleaseDate string `json:"releaseDate"`
	Text        string `json:"text"`
	Link        string `json:"link"`
}

// Добавление песни

type AddSong struct {
	Group string `json:"group"`
	Song  string `json:"song"`
}
