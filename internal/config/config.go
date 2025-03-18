package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const configFileName = ".gatorconfig.json"

func getConfigFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	filePath := filepath.Join(home, configFileName)
	return filePath, nil
}

func writeConfig(configuration Config) error {

	fileName, err := getConfigFilePath()

	data, err := json.Marshal(configuration)
	if err != nil {
		return err
	}
	fmt.Println("Config file:", fileName)
	err = os.WriteFile(fileName, data, 0777)
	if err != nil {
		return err
	}
	return nil
}

type Config struct {
	DBUrl           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func Read() (Config, error) {
	// Looks for config in home directory
	configFileLocation, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}
	configFile, err := os.ReadFile(configFileLocation)
	if err != nil {
		return Config{}, err
	}

	// Unmarshalling the data
	var config Config
	err = json.Unmarshal(configFile, &config)
	return config, nil
}

func SetUser(userName string, config *Config) error {
	(*config).CurrentUserName = userName
	err := writeConfig(*config)
	if err != nil {
		return err
	}
	return nil
}
