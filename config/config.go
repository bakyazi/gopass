package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

var once sync.Once
var config *Configuration

type Configuration struct {
	CredentialFile string
	SheetId        string
}

func (c *Configuration) Load(path string) {
	var err error
	if path == "" {
		path, err = os.UserHomeDir()
		if err != nil {
			fmt.Println("Cannot get user home directory", err.Error())
			os.Exit(1)
		}
		path = filepath.Join(path, ".gopass")
	}
	credFile := filepath.Join(path, "credentials.json")
	confFile := filepath.Join(path, "conf.json")
	if _, err := os.Stat(credFile); errors.Is(err, os.ErrNotExist) {
		fmt.Printf("Credential file not exists! %s\n", credFile)
		os.Exit(1)
	}

	if _, err := os.Stat(confFile); errors.Is(err, os.ErrNotExist) {
		fmt.Printf("Config file not exists! %s\n", confFile)
		os.Exit(1)
	}

	type confYml struct {
		SheetId string `json:"sheetId"`
	}
	var cy confYml
	data, err := os.ReadFile(confFile)
	if err != nil {
		fmt.Println("Cannot read config file", confFile, err.Error())
		os.Exit(1)
	}
	err = json.Unmarshal(data, &cy)
	if err != nil {
		fmt.Println("Cannot parse config file", confFile, err.Error())
		os.Exit(1)
	}
	c.SheetId = cy.SheetId
	c.CredentialFile = credFile

}

func GetConfig() *Configuration {
	once.Do(func() {
		config = &Configuration{}
		config.Load("")
	})
	return config
}
