package models

import (
	"database/sql"
)

// Source represents the source entity
type Source struct {
	ID        int
	Name      string
	Campaigns []*Campaign
}

func (s *Source) AddCampaign(campaign *Campaign, wasRelationSet bool) {
	if !wasRelationSet {
		campaign.AddSource(s, true)
	}

	s.Campaigns = append(s.Campaigns, campaign)
}

// SourceRepository represents a repository for saving sources in db
type SourceRepository struct {
	Db *sql.DB
}

// NewSourceFactory creates a new instance of SourceRepository
func NewSourceFactory(db *sql.DB) *SourceRepository {
	return &SourceRepository{db}
}

// Persist saves a single source in the database
func (sf *SourceRepository) Persist(source Source) error {
	_, err := sf.Db.Exec("INSERT INTO sources (name) VALUES (?)", source.Name)

	return err
}

// PersistAll performs a one query insert for an array of sources
func (sf *SourceRepository) PersistAll(sources []*Source) error {

	query := "INSERT INTO sources (name) VALUES"
	values := make([]interface{}, len(sources), len(sources))

	for i := 0; i < len(sources); i++ {
		if i > 0 {
			query += ","
		}
		query += "(?)"
		values[i] = sources[i].Name
	}

	_, err := sf.Db.Exec(query, values...)

	// yeah probably incorrect approach, should be rethought
	nextId, err := getNextAutoIncrementValue(sf.Db, "sources")

	for i := len(sources) - 1; i >= 0; i-- {
		nextId--
		sources[i].ID = int(nextId)
	}

	return err
}

// even if the insert query fails because of the unique constraint or because of the foreign key constraint, i think it is not a  problem because this operation seems ok to be idempotent.
//
// otherwise we would have to perform an existence query for each campaign, to check if it exists, and if not to first save the campaign entity itself, this seems not performant to me
//
// but firstly we insert in db all campaigns and all sources, and only after we call this method. So its safe to do a big insert like this.
// By doing this way, we would not have to do multiple separate insert queries for each source -> campaign relationship

// Saves the relationship of the source to campaigns in the db
// Ignores inserts that fail because of the reasons mentioned above
// Performs a big insert with all the campaigns for a source in the sources_associated_campaigns
func (sf *SourceRepository) PersistCampaignsRelation(source Source) error {
	query := "INSERT IGNORE INTO sources_associated_campaigns (campaign_id, source_id) VALUES"
	values := make([]interface{}, 0, 2*len(source.Campaigns))

	for i := 0; i < len(source.Campaigns); i++ {
		if i > 0 {
			query += ","
		}
		query += "(?,?)"

		values = append(values, source.Campaigns[i].ID)
		values = append(values, source.ID)
	}

	_, err := sf.Db.Exec(query, values...)

	return err
}
