package cli

import (
	"encoding/json"
	"math/big"
	"os"
	"path/filepath"

	"github.com/kirsle/configdir"
	"github.com/tmc/keyring"
)

var (
	// ConfigID is the general-purpose name we use to identify ourselves.
	// Can be overridden on an per-application basis from Init(), if
	// desired.
	ConfigID = "go-keepassrpc"
)

// Configuration represents our saved username and session key state.
type Configuration struct {
	Username string

	sessionKey *big.Int
	file       string
}

// Save checkpoints our configuration to disk, typically called after a
// successful SRP negotiation.
func (config *Configuration) Save() error {
	f, err := os.Create(config.file)
	if err != nil {
		return err
	}
	defer f.Close()
	encoder := json.NewEncoder(f)
	encoder.Encode(config)
	return keyring.Set(ConfigID, config.Username, config.sessionKey.Text(16))
}

// LoadConfig reads both our on-disk configuration state as well as the
// secret stored in the keyring.
func LoadConfig() (*Configuration, error) {
	var config Configuration

	configPath := configdir.LocalConfig(ConfigID)
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
		sessionKey, err := keyring.Get(ConfigID, config.Username)
		if err != nil {
			return nil, err
		}
		config.sessionKey, _ = new(big.Int).SetString(sessionKey, 16)
	}
	return &config, nil
}
