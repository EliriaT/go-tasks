package api

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"math/rand"
	"net/http/httptest"
	"testing"
)

var app = fiber.New()

// Is it all right if in the benchmark test  I am just making a request to the already running server?😄
func BenchmarkGetCampaignsBySource(b *testing.B) {
	randomSourceId := rand.Intn(100) + 1
	requestURL := fmt.Sprintf("http://localhost:8080/sources/%d", randomSourceId)
	req := httptest.NewRequest("GET", requestURL, nil)

	for i := 0; i < b.N; i++ {
		app.Test(req)
	}
}
