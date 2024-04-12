package models

import (
	"database/sql"
)

// Campaign represents the source entity
type Campaign struct {
	ID      int
	Name    string
	Sources []*Source `json:"-"`
}

func (c *Campaign) AddSource(source *Source, wasRelationSet bool) {
	if !wasRelationSet {
		source.AddCampaign(c, true)
	}

	c.Sources = append(c.Sources, source)
}

// SourceRepository represents a repository for saving campaigns in db
type CampaignRepository struct {
	Db *sql.DB
}

// NewCampaignFactory creates a new instance of CampaignRepository
func NewCampaignFactory(db *sql.DB) *CampaignRepository {
	return &CampaignRepository{db}
}

// Persist saves a single campaign in the database
func (cf *CampaignRepository) Persist(campaign Campaign) error {
	_, err := cf.Db.Exec("INSERT INTO campaigns (name) VALUES (?)", campaign.Name)

	return err
}

// PersistAll performs a one query insert for an array of campaigns
func (cf *CampaignRepository) PersistAll(campaigns []*Campaign) error {
	query := "INSERT INTO campaigns (name) VALUES"
	values := make([]interface{}, len(campaigns), len(campaigns))

	for i := 0; i < len(campaigns); i++ {
		if i > 0 {
			query += ","
		}
		query += "(?)"
		values[i] = campaigns[i].Name
	}

	_, err := cf.Db.Exec(query, values...)

	// yeah probably incorrect approach, should be rethought
	nextId, err := getNextAutoIncrementValue(cf.Db, "campaigns")

	for i := len(campaigns) - 1; i >= 0; i-- {
		nextId--
		campaigns[i].ID = int(nextId)
	}

	return err
}

// even if the insert query fails because of the unique constraint or because of the foreign key constraint, i think it is not a  problem because this operation seems ok to be idempotent.
//
// otherwise we would have to perform an existence query for each source, to check if it exists, and if not to first save the source entity itself, and this seems not performant to me
//
// but firstly we insert in db all campaigns and all sources, and only after we call this method. So its safe to do a big insert like this.
// By doing this way, we would not have to do multiple separate insert queries for each source -> campaign relationship

// Saves the relationship of the source to campaigns in the db
// Ignores inserts that fail because of the reasons mentioned above
// Performs a big insert with all the sources for a campaign in the sources_associated_campaigns
func (cf *CampaignRepository) PersistSourcesRelation(campaign Campaign) error {
	query := "INSERT IGNORE INTO sources_associated_campaigns (campaign_id, source_id) VALUES"
	values := make([]interface{}, 0, 2*len(campaign.Sources))

	for i := 0; i < len(campaign.Sources); i++ {
		if i > 0 {
			query += ","
		}
		query += "(?,?)"

		values = append(values, campaign.ID)
		values = append(values, campaign.Sources[i].ID)
	}
	//log.Println(query)
	//log.Println(values)

	_, err := cf.Db.Exec(query, values...)

	return err
}
