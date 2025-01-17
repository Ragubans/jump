package config

import (
	"os"
	"os/user"
	"path/filepath"

	"github.com/gsamokovarov/jump/scoring"
)

const (
	defaultScoreFile  = "scores.json"
	defaultSearchFile = "search.json"
	defaultPinsFile   = "pins.json"
	defaultHomeDir    = ".jump"
	defaultXDGDir     = "jump"
)

// Config represents the config directory and all the miscellaneous
// configuration files we can have in there.
type Config interface {
	ReadEntries() (scoring.Entries, error)
	WriteEntries(scoring.Entries) error

	ReadSearch() Search
	WriteSearch(string, int) error

	ReadPins() (map[string]string, error)
	FindPin(string) (string, bool)
	WritePin(string, string) error
	RemovePin(string) error
}

type fileConfig struct {
	Dir    string
	Scores string
	Search string
	Pins   string
}

// Setup setups the config folder from a directory path.
//
// If the directories don't already exists, they are created and if the score
// file is present, it is loaded.
func Setup(dir string) (Config, error) {
	// We get the directory check for free form os.MkdirAll.
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}

	scores := filepath.Join(dir, defaultScoreFile)
	search := filepath.Join(dir, defaultSearchFile)
	pins := filepath.Join(dir, defaultPinsFile)

	return &fileConfig{dir, scores, search, pins}, nil
}

// SetupDefault setups the config folder from a directory path.
//
// If the directory path is an empty string, the path is automatically guessed.
func SetupDefault(dir string) (Config, error) {
	dir, err := findConfigDir(dir)
	if err != nil {
		return nil, err
	}

	return Setup(dir)
}

// findConfigDir finds the jump configuration directory.
//
// The search algorithm tries the directories in order:
//
// - $JUMP_HOME (if given)
// - $HOME/.jump (if already exists)
// - $XDG_CONFIG_HOME/jump (prefer for new installs)
//
// We're moving towards XDG, but for existing installs or non-XDG supported
// systems, the ~/.jump dur will be used.
func findConfigDir(dir string) (string, error) {
	if dir != "" {
		return dir, nil
	}

	home, err := homeDir()
	if err != nil {
		return dir, err
	}

	configDir := filepath.Join(home, defaultHomeDir)
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		if xdgHome := os.Getenv("XDG_CONFIG_HOME"); xdgHome != "" {
			configDir = filepath.Join(xdgHome, defaultXDGDir)
		}
	}

	return configDir, nil
}

// See https://github.com/golang/go/issues/26463
func homeDir() (string, error) {
	home := os.Getenv("HOME")
	if home != "" {
		return home, nil
	}
	usr, err := user.Current()
	if err != nil {
		return "", err
	}

	return usr.HomeDir, nil
}
