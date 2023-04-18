package config

import (
	"strings"
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

type ProjectsHashMap map[string]*Project

type Config struct {
	Projects ProjectsHashMap         `json:"projects"`
	Teams    map[string]*CachedTeams `json:"teams"`
}

func (c *Config) GetProject(wd string) *Project {
	k := strings.ToLower(wd)

	return c.Projects[k]
}

func (c *Config) SetProject(wd string, project *Project) {
	k := strings.ToLower(wd)

	c.Projects[k] = project
}

func (c *Config) Update() {
	viper.Set("projects", c.Projects)
	viper.Set("teams", c.Teams)
}

type Project struct {
	TeamID  string `json:"teamId"`
	WorkDir string `json:"workdir"`
	ApiKey  string `json:"apiKey"`
}

func LoadDefaults() (Config, error) {
	config := Config{
		Projects: make(map[string]*Project),
		Teams:    make(map[string]*CachedTeams),
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
