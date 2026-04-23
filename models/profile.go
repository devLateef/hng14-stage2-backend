package models

import (
	"time"

	"github.com/google/uuid"
)

type Profile struct {
	ID                 uuid.UUID `gorm:"type:uuid;primaryKey"                          json:"id"`
	Name               string    `gorm:"uniqueIndex;not null"                          json:"name"`
	Gender             string    `gorm:"not null;index"                                json:"gender"`
	GenderProbability  float64   `gorm:"not null"                                      json:"gender_probability"`
	Age                int       `gorm:"not null;index"                                json:"age"`
	AgeGroup           string    `gorm:"not null;index"                                json:"age_group"`
	CountryID          string    `gorm:"type:varchar(2);not null;index"                json:"country_id"`
	CountryName        string    `gorm:"not null"                                      json:"country_name"`
	CountryProbability float64   `gorm:"not null"                                      json:"country_probability"`
	CreatedAt          time.Time `gorm:"autoCreateTime;index"                          json:"created_at"`
}
