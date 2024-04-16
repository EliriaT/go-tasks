package seed

import (
	"github.com/EliriaT/go-tasks/db/models"
	"math/rand"
)

// SourceSeeder helps with getting a random source
type SourceSeeder struct {
	Seeder
}

// GetNSources generates a slice of N random source names
func (s *SourceSeeder) GetNSources(n int) []*models.Source {
	var sources []*models.Source

	for i := 0; i < n; i++ {
		nameLength := rand.Intn(25) + 1
		sources = append(sources, &models.Source{Name: s.GetRandomName(nameLength)}) // Adjust the length as needed
	}
	return sources
}
