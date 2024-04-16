package main

import (
	"database/sql"
	"fmt"
	"github.com/EliriaT/go-tasks/api"
	"github.com/EliriaT/go-tasks/db/models"
	"github.com/EliriaT/go-tasks/db/seed"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
	"log"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	username := "user"
	password := "password"
	dbname := "sources"
	dbHost := "sources_db"
	dbPort := 3306
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", username, password, dbHost, dbPort, dbname)

	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		log.Fatalf("Error connecting to the database %v", err)
	}

	seedDb(db)

	server := api.NewServer(db)

	app := fiber.New()

	app.Get("/sources/:id<int>", server.GetCampaignsBySource)

	log.Fatal(app.Listen(":8080"))
}

func seedDb(db *sql.DB) {
	seeder := seed.Seeder{}
	sourceRepository := models.NewSourceRepository(db)
	sourceSeeder := seed.SourceSeeder{
		seeder,
	}

	campaignRepository := models.NewCampaignRepository(db)
	campaignSeeder := seed.CampaignSeeder{
		seeder,
	}

	domainRepository := models.NewDomainRepository(db)
	domainSeeder := seed.DomainSeeder{
		seeder,
	}

	campaigns := campaignSeeder.GetNCampaign(100)
	sources := sourceSeeder.GetNSources(100)
	domains := domainSeeder.GetNDomains(15)

	for i, _ := range campaigns {
		numberOfSources := rand.Intn(10)
		for j := 0; j < numberOfSources; j++ {
			randomSourceIndex := rand.Intn(100)
			campaigns[i].AddSource(sources[randomSourceIndex], false)
		}

		numberOfDomains := rand.Intn(10)
		for j := 0; j < numberOfDomains; j++ {
			randomDomainIndex := rand.Intn(15)
			campaigns[i].AddDomain(domains[randomDomainIndex])
		}
	}

	err := campaignRepository.PersistAll(campaigns)
	log.Println(err)

	err = sourceRepository.PersistAll(sources)
	log.Println(err)

	err = domainRepository.PersistAll(domains)
	log.Println(err)

	for _, campaign := range campaigns {
		if len(campaign.Sources) > 0 {
			err = campaignRepository.PersistSourcesRelation(*campaign)
			log.Println(err)
		}

		if len(campaign.Domains) > 0 {
			err = campaignRepository.PersistDomainsRelation(*campaign)
			log.Println(err)
		}
	}
}
