package api

import (
	"github.com/EliriaT/go-tasks/db/models"
	"github.com/EliriaT/go-tasks/db/seed"
	"github.com/gofiber/fiber/v2"
	"math/rand"
	"testing"
)

var app = fiber.New()

// Is it all right if in the benchmark test  I am just making a request to the already running server?😄
//func BenchmarkGetCampaignsBySource(b *testing.B) {
//	randomSourceId := rand.Intn(100) + 1
//	requestURL := fmt.Sprintf("http://localhost:8080/sources/%d", randomSourceId)
//	req := httptest.NewRequest("GET", requestURL, nil)
//
//	for i := 0; i < b.N; i++ {
//		app.Test(req)
//	}
//}

func BenchmarkFilterCampaignsViaMap(b *testing.B) {
	campaigns := getCampaigns()

	for i := 0; i < b.N; i++ {
		filterCampaignsViaMap("google.com", campaigns)
	}
}

func BenchmarkFilterCampaignsViaSlice(b *testing.B) {
	campaigns := getCampaigns()

	for i := 0; i < b.N; i++ {
		filterCampaignsViaSlice("google.com", campaigns)
	}
}

func getCampaigns() []*models.Campaign {
	seeder := seed.Seeder{}
	campaignSeeder := seed.CampaignSeeder{
		seeder,
	}
	domainSeeder := seed.DomainSeeder{
		seeder,
	}

	campaigns := campaignSeeder.GetNCampaign(1000)
	domains := domainSeeder.GetNDomains(15)

	for i, _ := range campaigns {
		numberOfDomains := rand.Intn(10)

		for j := 0; j < numberOfDomains; j++ {
			randomDomainIndex := rand.Intn(15)
			campaigns[i].AddDomain(domains[randomDomainIndex])
		}
	}

	return campaigns
}
