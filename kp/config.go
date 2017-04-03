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
	"github.com/logic/go-keepassrpc/keepassrpc"
)

type configuration struct {
	Username   string
	Password   string
	Value      *big.Int
	SessionKey *big.Int

	client *keepassrpc.Client
	file   string
}

func (config *configuration) prompt() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter the code provided by KeePass: ")
	text, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	config.Password = strings.TrimSpace(text)
	return config.Password, nil
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
	return nil
}

func loadConfig() (*configuration, error) {
	config := new(configuration)

	configPath := configdir.LocalConfig("kpcli")
	if err := configdir.MakePath(configPath); err != nil {
		panic(err)
	}

	config.file = filepath.Join(configPath, "settings.json")
	f, err := os.Open(config.file)
	if err != nil {
		if os.IsNotExist(err) {
			return config, nil
		}
		return nil, err
	}
	defer f.Close()
	d := json.NewDecoder(f)
	d.Decode(&config)
	return config, nil
}
