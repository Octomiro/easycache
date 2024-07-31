package easycache

import (
	"time"

	"github.com/patrickmn/go-cache"
)

type EasyCache struct {
	cache *cache.Cache
	// TTL for each cache entry
	TimeToLive int
	// Interval of cache evictions
	CleanUpInterval int
	// Http Status code limit for which responses are cached. Defaults to 300
	CacheIfStatusCodeLessThan int
	// List of endpoints (Paths) to ignore
	IgnoreEndpoints map[string]interface{}
	// To dis/able logging
	Logging bool
}

type CacheConfig struct {
	TimeToLive                int
	CleanUpInterval           int
	CacheIfStatusCodeLessThan int
	IgnoreEndpoints           map[string]interface{}
	Logging                   bool
}

func NewCache(conf CacheConfig) EasyCache {
	if conf.TimeToLive == 0 {
		conf.TimeToLive = 5
	}

	if conf.CleanUpInterval == 0 {
		conf.CleanUpInterval = 10
	}

	if conf.CacheIfStatusCodeLessThan == 0 {
		conf.CacheIfStatusCodeLessThan = 300
	}

	return EasyCache{
		cache: cache.New(time.Minute*time.Duration(conf.TimeToLive),
			time.Minute*time.Duration(conf.CleanUpInterval)),
		TimeToLive:                conf.TimeToLive,
		CleanUpInterval:           conf.CleanUpInterval,
		CacheIfStatusCodeLessThan: conf.CacheIfStatusCodeLessThan,
		IgnoreEndpoints:           conf.IgnoreEndpoints,
		Logging:                   conf.Logging,
	}
}

func (ec *EasyCache) Cache() *cache.Cache {
	return ec.cache
}
