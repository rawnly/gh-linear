package config

import (
	"time"

	"github.com/Rawnly/gh-linear/linear"
	"github.com/spf13/viper"
)

type CachedTeams struct {
	LastFetch int64         `json:"lastFetch"`
	Data      []linear.Team `json:"data"`
}

func (c *CachedTeams) IsExpired() bool {
	// now - lastFetch > 12h
	ttl := int64(60 * 60 * 12)

	return c.LastFetch+ttl < time.Now().Unix()
}

// Set updates the cached teams, does not save to disk
func (c *CachedTeams) Set(teams []linear.Team) {
	c.LastFetch = time.Now().Unix()
	c.Data = teams
}

type Config struct {
	Projects map[string]*Project `json:"projects"`
	Teams    CachedTeams         `json:"teams"`
}

func (c *Config) Update() {
	viper.Set("projects", c.Projects)
	viper.Set("teams", c.Teams)
}

type Project struct {
	TeamID  string `json:"teamId"`
	WorkDir string `json:"workdir"`
}

func LoadDefaults() (Config, error) {
	config := Config{
		Projects: make(map[string]*Project),
		Teams: CachedTeams{
			LastFetch: -1,
			Data:      nil,
		},
	}

	viper.SetDefault("teams", config.Teams)
	viper.SetDefault("projects", config.Projects)

	err := viper.WriteConfig()

	return config, err
}

func Load() (Config, error) {
	var C Config

	err := viper.Unmarshal(&C)

	return C, err
}

func (c *Config) Save() error {
	viper.Set("projects", c.Projects)
	return viper.WriteConfig()
}
