package handlers

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Client represents a client entity.
type Client struct {
	Id                int64
	ContractId        string
	PhoneNumber       string
	Address           string
	PaymentSum        string
	Comment           string
	Location          string
	LocationLatitude  string
	LocationLongitude string
	AddressFotoPath   string
	PaymentFotoPath   string
	UserId            int
}

// CacheItem represents an item stored in the cache.
type CacheItem struct {
	value      string
	expiration int64 // Unix timestamp (seconds)
}

// Cache represents an in-memory cache.
type Cache struct {
	items map[string]CacheItem
	mutex sync.RWMutex
}

// NewCache creates a new instance of Cache.
func NewCache() *Cache {
	return &Cache{
		items: make(map[string]CacheItem),
	}
}

// SetClientInfo sets a specific field of a client in the cache.
func (c *Cache) SetClientInfo(clientID int64, field string, value string, expiration int64) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	cacheKey := fmt.Sprintf("%d_%s", clientID, field)
	c.items[cacheKey] = CacheItem{
		value:      value,
		expiration: time.Now().Unix() + expiration,
	}
}

// GetClientInfo retrieves a specific field of a client from the cache.
func (c *Cache) GetClientInfo(clientID int64, field string) (string, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	cacheKey := fmt.Sprintf("%d_%s", clientID, field)
	item, found := c.items[cacheKey]
	if !found {
		return "", false
	}

	if item.expiration != 0 && time.Now().Unix() > item.expiration {
		// Item has expired
		delete(c.items, cacheKey)
		return "0", false
	}

	return item.value, true
}

// GetAllClientData retrieves all data for a specific client from the cache.
func (c *Cache) GetAllClientData(clientID int64) map[string]interface{} {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	clientData := make(map[string]interface{})
	for key, item := range c.items {
		parts := strings.Split(key, "_")
		if len(parts) == 2 {
			cachedClientID, err := strconv.ParseInt(parts[0], 10, 64)
			if err != nil {
				// Handle error
				continue
			}
			if cachedClientID == clientID {
				clientData[parts[1]] = item.value
			}
		}
	}
	return clientData
}
