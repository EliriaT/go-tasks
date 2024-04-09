# Query 1
SELECT s.name, count(sac.campaign_id) AS nr_campaigns
FROM sources_associated_campaigns AS sac
RIGHT JOIN sources AS s ON s.id = sac.source_id
GROUP BY  s.id
ORDER BY nr_campaigns DESC
LIMIT 5;

# Query 2
SELECT c.name, c.id FROM campaigns AS c
                             LEFT JOIN sources_associated_campaigns AS sac ON c.id = sac.campaign_id
WHERE sac.source_id IS NULL;

# Query 3
SELECT s.name FROM sources AS s
UNION
SELECT c.name FROM campaigns AS c;
