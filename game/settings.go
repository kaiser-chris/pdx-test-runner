package game

import (
	"encoding/json"
	"fmt"
	"golang.org/x/sys/windows/registry"
	"os"
	"path/filepath"
	"strings"
)

type Type int

const ( // iota is reset to 0
	Victoria3 Type = iota
	CrusaderKings3
)

const windowsDocumentsEnv = "HOMEPATH"
const userDocumentsVariable = "%USER_DOCUMENTS%"

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

	if strings.Contains(launcherSettings.DataPath, userDocumentsVariable) {
		documentsFolderKey, err := registry.OpenKey(registry.CURRENT_USER, `SOFTWARE\Microsoft\Windows\CurrentVersion\Explorer\User Shell Folders`, registry.QUERY_VALUE)
		if err != nil {
			return nil, fmt.Errorf("could not query registry for user directory (%s): %v", userDocumentsVariable, err)
		}
		documentsFolder, _, err := documentsFolderKey.GetStringValue("Personal")
		if err != nil {
			return nil, fmt.Errorf("could not query registry for user directory (%s): %v", userDocumentsVariable, err)
		}
		launcherSettings.DataPath = strings.ReplaceAll(launcherSettings.DataPath, userDocumentsVariable, documentsFolder)
	}
	launcherSettings.DataPath = filepath.Clean(launcherSettings.DataPath)

	switch launcherSettings.GameId {
	case "victoria3":
		launcherSettings.Type = Victoria3
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
