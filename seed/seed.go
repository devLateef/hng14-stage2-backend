package seed

import (
	"encoding/json"
	"log"
	"os"
	"time"

	"insight-api/config"
	"insight-api/models"

	"github.com/google/uuid"
)

// seedFile matches the JSON structure: {"profiles": [...]}
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

// newUUIDv7 generates a UUID v7 (time-ordered).
// Falls back to uuid.New() (v4) if the runtime doesn't support v7.
func newUUIDv7() uuid.UUID {
	id, err := uuid.NewV7()
	if err != nil {
		return uuid.New()
	}
	return id
}

// SeedProfiles reads profiles from filePath and inserts them into the DB.
// Re-running is safe — existing names are skipped.
func SeedProfiles(filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	var sf seedFile
	if err := json.Unmarshal(data, &sf); err != nil {
		return err
	}

	inserted := 0
	skipped := 0

	for _, sp := range sf.Profiles {
		var existing models.Profile
		result := config.DB.Where("name = ?", sp.Name).First(&existing)
		if result.Error == nil {
			skipped++
			continue // already exists
		}

		profile := models.Profile{
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
		}

		if err := config.DB.Create(&profile).Error; err != nil {
			log.Printf("Failed to insert profile %q: %v", sp.Name, err)
			continue
		}
		inserted++
	}

	log.Printf("Seed complete: %d inserted, %d skipped (duplicates)", inserted, skipped)
	return nil
}
