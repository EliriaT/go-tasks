package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/EliriaT/go-tasks/db/models"
	"github.com/gofiber/fiber/v2"
	"log"
	"strconv"
	"strings"
)

type Server struct {
	cache      *Cache
	sourceRepo models.SourceRepository
}

func NewServer(db *sql.DB) Server {
	return Server{
		cache:      NewCache(),
		sourceRepo: models.NewSourceRepository(db),
	}
}

func (s Server) GetCampaignsBySource(c *fiber.Ctx) error {
	sourceId := c.Params("id")
	domainWhitelist := strings.ToLower(c.Query("domain"))

	sourceIdInt, err := strconv.Atoi(sourceId)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Source id is not an  int")
	}

	campaigns, present := s.cache.Get(sourceIdInt)

	if !present {
		source, err := s.sourceRepo.GetSourceWithCampaignsAndDomains(sourceIdInt)
		if err != nil {
			if errors.Is(err, models.ErrNotFound) {
				return c.Status(fiber.StatusNotFound).SendString("No such source exist")
			}

			log.Println(err)
			return c.Status(fiber.StatusInternalServerError).SendString("Something went wrong")
		}

		s.cache.Put(source)
		campaigns = source.Campaigns
	}

	campaigns = filterCampaignsByWhitelist(domainWhitelist, campaigns)

	jsonData, err := json.Marshal(campaigns)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Something went wrong")
	}

	c.Set("Content-Type", "application/json")

	return c.Send(jsonData)
}

func filterCampaignsByWhitelist(domain string, campaigns []*models.Campaign) []*models.Campaign {
	if domain == "" {
		return campaigns
	}

	filteredCampaigns := make([]*models.Campaign, 0, len(campaigns))

	for _, campaign := range campaigns {
		if passesBlackFilter(campaign, domain) || passesWhiteFilter(campaign, domain) {
			filteredCampaigns = append(filteredCampaigns, campaign)
		}
	}
	return filteredCampaigns
}

func passesWhiteFilter(campaign *models.Campaign, domain string) bool {
	if campaign.ListType == models.WhiteList {
		if _, ok := campaign.List[domain]; ok {
			return true
		}
	}

	return false
}

func passesBlackFilter(campaign *models.Campaign, domain string) bool {
	if campaign.ListType == models.BlackList {
		if _, ok := campaign.List[domain]; !ok {
			return true
		}
	}

	return false
}
