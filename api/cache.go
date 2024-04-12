package api

import (
	"github.com/EliriaT/go-tasks/db/models"
	"sync"
)

type Cache struct {
	mu                 sync.RWMutex
	campaignsBySources map[int][]*models.Campaign
}

func NewCache() *Cache {
	return &Cache{
		campaignsBySources: make(map[int][]*models.Campaign),
	}
}

func (c *Cache) Put(source models.Source) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.campaignsBySources[source.ID] = source.Campaigns
}

func (c *Cache) Get(id int) ([]*models.Campaign, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	campaigns, ok := c.campaignsBySources[id]

	return campaigns, ok
}
