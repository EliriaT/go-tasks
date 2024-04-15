CREATE TABLE campaigns_associated_domains (
          campaign_id INT NOT NULL ,
          domain_id INT NOT NULL ,
          FOREIGN KEY (campaign_id) REFERENCES campaigns(id) on delete cascade,
          FOREIGN KEY (domain_id) REFERENCES domains(id) on delete cascade,

          CONSTRAINT unique_campaign_domain UNIQUE (campaign_id, domain_id)
);