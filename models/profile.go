package models

import (
	"time"

	"github.com/google/uuid"
)

type Profile struct {
	ID                 uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Name               string    `gorm:"unique;not null" json:"name"`
	Gender             string    `json:"gender"`
	GenderProbability  float64   `json:"gender_probability"`
	Age                int       `json:"age"`
	AgeGroup           string    `json:"age_group"`
	CountryID          string    `json:"country_id"`
	CountryName        string    `json:"country_name"`
	CountryProbability float64   `json:"country_probability"`
	CreatedAt          time.Time `gorm:"autoCreateTime" json:"created_at"`
}
