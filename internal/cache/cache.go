package cache

import (
	"sync"
	"time"
)

type TokenCacheItem struct {
	expiresAt time.Time
	token     string
}

type TokenCache struct {
	tokens map[string]TokenCacheItem
	mu     sync.RWMutex
}

func NewTokenCache() *TokenCache {
	return &TokenCache{tokens: make(map[string]TokenCacheItem)}
}

func (c *TokenCache) Set(key, token string, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.tokens[key] = TokenCacheItem{
		token:     token,
		expiresAt: time.Now().Add(ttl),
	}
}

func (c *TokenCache) Get(key string) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	item, ok := c.tokens[key]
	if !ok || time.Now().After(item.expiresAt) {
		return "", false
	}
	return item.token, true
}

func (c *TokenCache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.tokens, key)
}

func (c *TokenCache) Cleanup(interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		for range ticker.C {
			c.mu.Lock()
			now := time.Now()
			for k, v := range c.tokens {
				if now.After(v.expiresAt) {
					delete(c.tokens, k)
				}
			}
			c.mu.Unlock()
		}
	}()
}
