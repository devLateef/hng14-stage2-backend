package seed

import (
	"encoding/json"
	"os"

	"insight-api/config"
	"insight-api/models"

	"github.com/google/uuid"
)

type seedFile struct {
	Profiles []models.Profile `json:"profiles"`
}

func SeedProfiles(filePath string) error {
	file, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	var data seedFile
	if err := json.Unmarshal(file, &data); err != nil {
		return err
	}

	for _, p := range data.Profiles {
		var existing models.Profile
		if err := config.DB.Where("name = ?", p.Name).First(&existing).Error; err == nil {
			// Record already exists, skip
			continue
		}

		// Generate UUID v7
		p.ID, err = uuid.NewV7()
		if err != nil {
			// Fall back to v4 if v7 generation fails
			p.ID = uuid.New()
		}

		config.DB.Create(&p)
	}

	return nil
}
