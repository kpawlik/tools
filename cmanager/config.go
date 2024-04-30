package cmanager

import (
	"encoding/json"
	"fmt"
	"os"
)

type User struct {
	UserName     string `json:"userName"`
	ImageName    string `json:"imageName"`
	ImageTag     string `json:"imageTag"`
	HostName     string `json:"hostName"`
	MemcachePort string `json:"memcachePort"`
	ApachePort   string `json:"apachePort"`
	PgDatabase   string `json:"pgDatabase"`
	PgUser       string `json:"pgUser"`
	PgHost       string `json:"pgHost"`
	PgPort       string `json:"pgPort"`
	PgPassword   string `json:"PgPassword"`
}

type Config struct {
	ImageName  string
	ImageTag   string
	HostName   string `json:"hostName"`
	UsersDir   string `json:"userDir"`
	RepoName   string `json:"repoName"`
	ComposeDir string `json:"composeDir"`
	PgDatabase string `json:"pgDatabase"`
	PgUser     string `json:"pgUser"`
	PgHost     string `json:"pgHost"`
	PgPort     string `json:"pgPort"`
	PgPassword string `json:"PgPassword"`
	Users      []User `json:"users"`
}

func LoadConfig(fileName string) (cfg Config, err error) {
	var (
		buff []byte
	)
	if buff, err = os.ReadFile(fileName); err != nil {
		err = fmt.Errorf("error reading file %s %v", fileName, err)
		return
	}
	if err = json.Unmarshal(buff, &cfg); err != nil {
		err = fmt.Errorf("error unmarshal config from %s. %v", fileName, err)
		return
	}
	return
}
