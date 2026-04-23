package seed

import (
	_ "embed"
	"encoding/json"
	"log"
	"time"

	"insight-api/config"
	"insight-api/models"

	"github.com/google/uuid"
	"gorm.io/gorm/clause"
)

//go:embed profiles.json
var profilesData []byte

type seedFile struct {
	Profiles []seedProfile `json:"profiles"`
}

type seedProfile struct {
	Name               string  `json:"name"`
	Gender             string  `json:"gender"`
	GenderProbability  float64 `json:"gender_probability"`
	Age                int     `json:"age"`
	AgeGroup           string  `json:"age_group"`
	CountryID          string  `json:"country_id"`
	CountryName        string  `json:"country_name"`
	CountryProbability float64 `json:"country_probability"`
}

func newUUIDv7() uuid.UUID {
	id, err := uuid.NewV7()
	if err != nil {
		return uuid.New()
	}
	return id
}

// SeedProfiles bulk-inserts all profiles using ON CONFLICT (name) DO NOTHING.
// Fast on first run, instant on subsequent runs.
func SeedProfiles() error {
	var sf seedFile
	if err := json.Unmarshal(profilesData, &sf); err != nil {
		return err
	}

	if len(sf.Profiles) == 0 {
		return nil
	}

	profiles := make([]models.Profile, 0, len(sf.Profiles))
	for _, sp := range sf.Profiles {
		profiles = append(profiles, models.Profile{
			ID:                 newUUIDv7(),
			Name:               sp.Name,
			Gender:             sp.Gender,
			GenderProbability:  sp.GenderProbability,
			Age:                sp.Age,
			AgeGroup:           sp.AgeGroup,
			CountryID:          sp.CountryID,
			CountryName:        sp.CountryName,
			CountryProbability: sp.CountryProbability,
			CreatedAt:          time.Now().UTC(),
		})
	}

	// Bulk insert in batches of 500 rows.
	// ON CONFLICT (name) DO NOTHING skips duplicates without error.
	result := config.DB.
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "name"}},
			DoNothing: true,
		}).
		CreateInBatches(profiles, 500)

	if result.Error != nil {
		return result.Error
	}

	log.Printf("Seed complete: %d profiles processed", len(profiles))
	return nil
}
