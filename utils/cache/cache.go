package cache

import (
	"time"

	"github.com/Rawnly/gh-linear/utils"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Cache struct {
	Key       string      `json:"key"`
	Value     interface{} `json:"value"`
	ExpiresAt int64       `json:"expires_at"`
}

func (c *Cache) IsExpired() bool {
	return c.ExpiresAt < time.Now().Unix()
}

func (c *Cache) SetExpiration(duration time.Duration) {
	c.ExpiresAt = time.Now().Add(duration).Unix()
}

func (c *Cache) SetExpirationFromNow(duration time.Duration) {
	c.ExpiresAt = time.Now().Add(duration).Unix()
}

func (c *Cache) SetExpirationFromNowInMinutes(seconds int64) {
	c.SetExpirationFromNow(time.Duration(seconds) * time.Minute)
}

func (c *Cache) SetExpirationFromNowInSeconds(seconds int64) {
	c.SetExpirationFromNow(time.Duration(seconds) * time.Second)
}

func (c *Cache) Mutate(value interface{}, ttl int64) error {
	c.Value = value
	c.SetExpirationFromNowInMinutes(ttl)

	cache := getCache()
	cache[c.Key] = *c

	viper.Set("cache", cache)

	return viper.WriteConfig()
}

func Read(key string, default_ttl int64) (Cache, error) {
	cache := getCache()
	cached := cache[key]
	cached.Key = key

	logrus.Debug("[[cache hit]]")
	logrus.Debug(utils.PrettyJSON(cached))

	return cached, nil
}

func New(key string, value interface{}, ttl int64) Cache {
	return Cache{
		Key:       key,
		Value:     value,
		ExpiresAt: time.Now().Add(time.Duration(ttl) * time.Minute).Unix(),
	}
}

func getCache() map[string]Cache {
	cache := viper.Get("cache")
	logrus.Debug(cache)

	if cache == nil {
		return make(map[string]Cache)
	}

	return cache.(map[string]Cache)
}
