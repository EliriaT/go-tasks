package seed

import (
	"github.com/EliriaT/go-tasks/db/models"
	"math/rand"
)

// CampaignSeeder helps with getting a random source
type CampaignSeeder struct {
	Seeder
	Repository models.CampaignRepository
}

// GetNSources generates a slice of N random source names
func (c *CampaignSeeder) GetNCampaign(n int) []*models.Campaign {
	var campaigns []*models.Campaign

	for i := 0; i < n; i++ {
		nameLength := rand.Intn(250) + 1
		campaigns = append(campaigns, &models.Campaign{Name: c.GetRandomName(nameLength)}) // Adjust the length as needed
	}

	return campaigns
}

func (c *CampaignSeeder) SeedInDb(campaigns []*models.Campaign) error {
	return c.Repository.PersistAll(campaigns)
}
