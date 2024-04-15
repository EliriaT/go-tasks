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

SELECT s.id, c.id, c.name, c.list_type, d.domain FROM  sources_associated_campaigns as sac
RIGHT JOIN sources as s ON s.id = sac.source_id
LEFT JOIN campaigns as c on c.id = sac.campaign_id
LEFT JOIN campaigns_associated_domains cad on c.id = cad.campaign_id
LEFT JOIN domains d on d.id = cad.domain_id
WHERE s.id = 1
ORDER BY c.name