package cache

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Cache handles local caching of Azure resource data
type Cache struct {
	dir string
}

// NewCache creates a new cache instance
func NewCache(dir string) (*Cache, error) {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create cache directory: %w", err)
	}

	return &Cache{dir: dir}, nil
}

// Close closes the cache (placeholder for future DB-based cache)
func (c *Cache) Close() error {
	return nil
}

// Get retrieves cached data
func (c *Cache) Get(key string, dest interface{}) error {
	path := filepath.Join(c.dir, key+".json")
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, dest)
}

// Set stores data in cache
func (c *Cache) Set(key string, data interface{}) error {
	path := filepath.Join(c.dir, key+".json")

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, jsonData, 0644)
}

// Has checks if key exists in cache
func (c *Cache) Has(key string) bool {
	path := filepath.Join(c.dir, key+".json")
	_, err := os.Stat(path)
	return err == nil
}

// GetWithTTL retrieves cached data with TTL check
func (c *Cache) GetWithTTL(key string, ttl time.Duration, dest interface{}) error {
	path := filepath.Join(c.dir, key+".json")

	info, err := os.Stat(path)
	if err != nil {
		return err
	}

	// Check if cache is stale
	if time.Since(info.ModTime()) > ttl {
		return fmt.Errorf("cache expired")
	}

	return c.Get(key, dest)
}

// Clear removes all cached data
func (c *Cache) Clear() error {
	return os.RemoveAll(c.dir)
}
