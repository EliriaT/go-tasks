package models

import "database/sql"

// Domain represents the domain entity
type Domain struct {
	ID        int
	Name      string
	Campaigns []*Campaign
}

func (d *Domain) AddCampaign(campaign *Campaign) {
	if campaign.ListType == WhiteList {
		d.Campaigns = append(d.Campaigns, campaign)
	}
}

// DomainRepository represents a repository for saving domains in db
type DomainRepository struct {
	Db *sql.DB
}

// NewDomainRepository creates a new instance of DomainRepository
func NewDomainRepository(db *sql.DB) DomainRepository {
	return DomainRepository{db}
}

// PersistAll performs a one query insert for an array of domains
func (cf *DomainRepository) PersistAll(domains []*Domain) error {
	query := "INSERT INTO domains (domain) VALUES"
	values := make([]interface{}, len(domains), len(domains))

	for i := 0; i < len(domains); i++ {
		if i > 0 {
			query += ","
		}
		query += "(?)"
		values[i] = domains[i].Name
	}

	_, err := cf.Db.Exec(query, values...)

	nextId, err := getNextAutoIncrementValue(cf.Db, "domains")

	for i := len(domains) - 1; i >= 0; i-- {
		nextId--
		domains[i].ID = int(nextId)
	}

	return err
}
