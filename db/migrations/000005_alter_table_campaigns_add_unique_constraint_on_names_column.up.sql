ALTER TABLE campaigns
    ADD CONSTRAINT unique_name UNIQUE (name);