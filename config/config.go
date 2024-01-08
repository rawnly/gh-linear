package config

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/Rawnly/gh-linear/linear"
	"github.com/spf13/viper"
)

type CachedTeams struct {
	ApiKey    string            `json:"key"`
	LastFetch int64             `json:"lastFetch"`
	Data      []linear.Team     `json:"data"`
	Workspace *linear.Workspace `json:"organization"`
}

func (c *CachedTeams) IsExpired() bool {
	// now - lastFetch > 12h
	ttl := int64(60 * 60 * 12)

	if c.LastFetch == 0 {
		return true
	}

	return c.LastFetch+ttl < time.Now().Unix()
}

// Set updates the cached teams, does not save to disk
func (c *CachedTeams) Set(key string, teams []linear.Team, workspace *linear.Workspace) {
	c.ApiKey = key
	c.LastFetch = time.Now().Unix()
	c.Data = teams
	c.Workspace = workspace
}

type ProjectsHashMap map[string]*Project

type Config struct {
	Projects ProjectsHashMap         `json:"projects"`
	Teams    map[string]*CachedTeams `json:"teams"`
}

func (c *Config) ResetRefresh() error {
	for _, t := range c.Teams {
		t.LastFetch = 0
	}

	c.Update()
	return c.Save()
}

func (c *Config) GetTeams() ([]CachedTeams, error) {
	hashmap := make(map[string]*CachedTeams)
	b, err := json.Marshal(viper.Get("teams"))
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(b, &hashmap); err != nil {
		return nil, err
	}

	var teams []CachedTeams

	for _, t := range c.Teams {
		teams = append(teams, *t)
	}

	return teams, nil
}

func (c *Config) GetProjects() []Project {
	var projects []Project

	for _, p := range c.Projects {
		projects = append(projects, *p)
	}

	return projects
}

func (c *Config) GetProject(wd string) *Project {
	k := strings.ToLower(wd)

	return c.Projects[k]
}

func (c *Config) SetProject(wd string, project *Project) {
	k := strings.ToLower(wd)

	c.Projects[k] = project
}

func (c *Config) GetTeam(teamId string) *CachedTeams {
	return c.Teams[teamId]
}

func (c *Config) SetTeam(teamId string, teams *CachedTeams) {
	c.Teams[teamId] = teams
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
	viper.Set("teams", c.Teams)

	return viper.WriteConfig()
}
