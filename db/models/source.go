package models

import (
	"database/sql"
	"errors"
)

var ErrNotFound = errors.New("Error Not Found")

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

// SourceRepository creates a new instance of SourceRepository
func NewSourceRepository(db *sql.DB) SourceRepository {
	return SourceRepository{db}
}

// we could also filter directly by domain:
// WHERE d.name = "google.com"
// and then store the result in cache and invalidate it when necessary
// But in case query by domain varies a lot (there are a lot of different domains that are queried), doing without `WHERE d.name = "google.com"` can be better I think
func (sf *SourceRepository) GetSourceWithCampaignsAndDomains(sourceId int) (Source, error) {
	query := `SELECT s.id, c.id, c.name, c.list_type, d.domain FROM  sources_associated_campaigns as sac
				RIGHT JOIN sources as s ON s.id = sac.source_id
				LEFT JOIN campaigns as c on c.id = sac.campaign_id
			 	LEFT JOIN campaigns_associated_domains cad on c.id = cad.campaign_id
			 	LEFT JOIN domains d on d.id = cad.domain_id
				WHERE s.id = ?
				ORDER BY c.name`

	rows, err := sf.Db.Query(query, sourceId)
	if err != nil {
		return Source{}, err
	}
	defer rows.Close()

	source := Source{}
	previousCampaignID := 0
	var campaign = Campaign{
		List: make(map[string]bool),
	}

	for rows.Next() {
		var domainColumn sql.NullString // the joined column can be null if the white/blacklist is empty
		var domain Domain

		if err := rows.Scan(&source.ID, &campaign.ID, &campaign.Name, &campaign.ListType, &domainColumn); err != nil {
			return Source{}, err
		}

		if domainColumn.Valid {
			domain.Name = domainColumn.String

			// this is because the query contains repeated campaign names and their different domains
			if campaign.ID == previousCampaignID {
				lastCampaignIndex := len(source.Campaigns) - 1
				source.Campaigns[lastCampaignIndex].AddDomain(&domain)
			}

			if campaign.ID != previousCampaignID {
				campaign.List = make(map[string]bool)
				campaign.AddDomain(&domain)
				copiedCampaign := campaign
				source.Campaigns = append(source.Campaigns, &copiedCampaign)
			}
		}

		previousCampaignID = campaign.ID
	}

	if source.ID == 0 {
		return Source{}, ErrNotFound
	}

	return source, err
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
