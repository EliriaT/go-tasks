package seed

import (
	"github.com/EliriaT/go-tasks/db/models"
	"math/rand"
)

// SourceSeeder helps with getting a random source
type SourceSeeder struct {
	Seeder
	Repository models.SourceRepository
}

// GetNSources generates a slice of N random source names
func (s *SourceSeeder) GetNSources(n int) []*models.Source {
	var sources []*models.Source

	for i := 0; i < n; i++ {
		nameLength := rand.Intn(250) + 1
		sources = append(sources, &models.Source{Name: s.GetRandomName(nameLength)}) // Adjust the length as needed
	}
	return sources
}

func (s *SourceSeeder) SeedInDb(sources []*models.Source) error {
	return s.Repository.PersistAll(sources)
}
