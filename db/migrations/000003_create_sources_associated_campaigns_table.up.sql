CREATE TABLE sources_associated_campaigns (
      campaign_id INT NOT NULL ,
      source_id INT NOT NULL ,
      FOREIGN KEY (campaign_id) REFERENCES campaigns(id) on delete cascade,
#     Should the associated campaigns be deleted when a source is deleted? Or should it be on delete restrict when deleting sources?
      FOREIGN KEY (source_id) REFERENCES sources(id),

      CONSTRAINT unique_campaign_source UNIQUE (campaign_id, source_id)
);