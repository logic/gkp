package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"strings"

	"github.com/kirsle/configdir"
	"github.com/tmc/keyring"
)

type configuration struct {
	Username string

	sessionKey *big.Int
	file       string
}

func (config *configuration) prompt() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter the code provided by KeePass: ")
	text, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(text), nil
}

// Save checkpoints our configuration to disk, typically called after a
// successful SRP negotiation.
func (config *configuration) Save() error {
	f, err := os.Create(config.file)
	if err != nil {
		return err
	}
	defer f.Close()
	encoder := json.NewEncoder(f)
	encoder.Encode(config)
	return keyring.Set("kp", config.Username, config.sessionKey.Text(16))
}

func loadConfig() (*configuration, error) {
	var config configuration

	configPath := configdir.LocalConfig("kp")
	if err := configdir.MakePath(configPath); err != nil {
		return nil, err
	}

	config.file = filepath.Join(configPath, "settings.json")
	f, err := os.Open(config.file)
	if err != nil {
		if os.IsNotExist(err) {
			return &config, nil
		}
		return nil, err
	}
	defer f.Close()
	d := json.NewDecoder(f)
	d.Decode(&config)
	if config.Username != "" {
		sessionKey, err := keyring.Get("kp", config.Username)
		if err != nil {
			return nil, err
		}
		config.sessionKey, _ = new(big.Int).SetString(sessionKey, 16)
	}
	return &config, nil
}
