package game

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Type int

const ( // iota is reset to 0
	Victoria3 Type = iota
	CrusaderKings3
)

type LauncherSettings struct {
	Type        Type
	GameId      string `json:"gameId"`
	DataPath    string `json:"gameDataPath"`
	ExecPath    string `json:"exePath"`
	ContentPath string `json:"dlcPath"`
}

func GetLauncherSettings(basePath string) (*LauncherSettings, error) {
	launcherDirectory := filepath.Join(basePath, "launcher")
	launcherSettingsPath := filepath.Join(launcherDirectory, "launcher-settings.json")
	var launcherSettings LauncherSettings
	content, err := os.ReadFile(launcherSettingsPath)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(content, &launcherSettings)
	if err != nil {
		return nil, err
	}

	switch launcherSettings.GameId {
	case "victoria3":
		launcherSettings.Type = Victoria3
		launcherSettings.DataPath = filepath.Clean(launcherSettings.DataPath)
		launcherSettings.ContentPath = filepath.Join(launcherDirectory, launcherSettings.ContentPath)
		launcherSettings.ExecPath = filepath.Join(launcherDirectory, launcherSettings.ExecPath)
	case "ck3":
		launcherSettings.Type = CrusaderKings3
		launcherSettings.DataPath = filepath.Clean(launcherSettings.DataPath)
		launcherSettings.ContentPath = filepath.Join(launcherDirectory, launcherSettings.ContentPath)
		launcherSettings.ExecPath = filepath.Join(launcherDirectory, launcherSettings.ExecPath)
	default:
		return nil, fmt.Errorf("unsupported game id: %s", launcherSettings.GameId)
	}

	return &launcherSettings, nil
}
