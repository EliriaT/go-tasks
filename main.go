package main

import (
	"database/sql"
	"fmt"
	"github.com/EliriaT/go-tasks/db/models"
	"github.com/EliriaT/go-tasks/db/seed"
	_ "github.com/go-sql-driver/mysql"
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

	sourceRepository := models.SourceRepository{
		db,
	}
	sourceSeeder := seed.SourceSeeder{
		seed.Seeder{},
		sourceRepository,
	}

	campaignRepository := models.CampaignRepository{
		db,
	}
	campaignSeeder := seed.CampaignSeeder{
		seed.Seeder{},
		campaignRepository,
	}

	campaigns := campaignSeeder.GetNCampaign(100)
	sources := sourceSeeder.GetNSources(100)

	for i, _ := range campaigns {
		numberOfSources := rand.Intn(10)

		for j := 0; j < numberOfSources; j++ {
			randomSourceIndex := rand.Intn(100)
			campaigns[i].AddSource(sources[randomSourceIndex], false)

		}
	}

	err = campaignRepository.PersistAll(campaigns)
	log.Println(err)

	err = sourceRepository.PersistAll(sources)
	log.Println(err)

	for _, campaign := range campaigns {
		if len(campaign.Sources) > 0 {
			err = campaignRepository.PersistSourcesRelation(*campaign)
			log.Println(err)
		}
	}
}
