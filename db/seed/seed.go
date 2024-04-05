package seed

import (
	"math/rand"
)

// Seeder  generates random data
type Seeder struct {
}

// GetRandomName generates a random string name of given length
func (s *Seeder) GetRandomName(length int) string {
	letters := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := make([]byte, length)

	for i := 0; i < length; i++ {
		bytes[i] = letters[rand.Intn(len(letters))]
	}

	return string(bytes)
}
