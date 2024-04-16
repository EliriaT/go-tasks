package seed

import (
	"github.com/EliriaT/go-tasks/db/models"
	"github.com/brianvoe/gofakeit/v6"
)

// DomainSeeder helps with getting a random domain
type DomainSeeder struct {
	Seeder
}

// GetNSources generates a slice of N random source names
func (d *DomainSeeder) GetNDomains(n int) []*models.Domain {
	var domains []*models.Domain

	for i := 0; i < n; i++ {
		domains = append(domains, &models.Domain{Name: gofakeit.DomainName()})
	}

	return domains
}
