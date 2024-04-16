package seed

import (
	"github.com/EliriaT/go-tasks/db/models"
	"math/rand"
)

// CampaignSeeder helps with getting a random source
type CampaignSeeder struct {
	Seeder
}

// GetNCampaign generates a slice of N random campaign names
func (c *CampaignSeeder) GetNCampaign(n int) []*models.Campaign {
	var campaigns []*models.Campaign

	for i := 0; i < n; i++ {
		nameLength := rand.Intn(25) + 1
		campaigns = append(campaigns, &models.Campaign{Name: c.GetRandomName(nameLength), ListType: models.ListType(rand.Intn(2)), List: make(map[string]bool)})
	}

	return campaigns
}
