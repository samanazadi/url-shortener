package models

type URL struct {
	ShortURL    string  `gorm:"primaryKey"`
	OriginalURL string  `gorm:"not null"`
	Visits      []Visit `gorm:"foreignKey:ShortURL;references:ShortURL"`
}
