package main

import (
	"bahmut.de/pdx-test-runner/logging"
	"golang.org/x/sys/windows/registry"
	"os"
	"path"
	"regexp"
)

func main() {
}

func getGamePath() string {
	//regex := regexp.MustCompile(`{\s+"path"\s*"(.*)".*\s+"label"\s+".*?"\s+"contentid"\s+".*?"\s+"totalsize"\s+".*?"\s+"update_clean_bytes_tally"\s+".*?"\s+"time_last_update_verified"\s+".*?"\s+"apps"\s+{(?s)[\s\"\d]*?"(545637013)"[\s"\d]*?}`)
	steamPathKey, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\WOW6432Node\Valve\Steam`, registry.QUERY_VALUE)
	if err != nil {
		logging.Fatal(err)
	}
	steamPath, _, err := steamPathKey.GetStringValue("InstallPath")
	if err != nil {
		logging.Fatal(err)
	}
	err = steamPathKey.Close()
	if err != nil {
		logging.Fatal(err)
	}
	libraryPath := path.Join(steamPath, "steamapps", "libraryfolders.vdf")
	libraryFile, err := os.Open(libraryPath)
	if err != nil {
		logging.Fatal(err)
	}
	vdfParser := vdf.NewParser(libraryFile)
	vdfContent, err := vdfParser.Parse()
	if err != nil {
		logging.Fatal(err)
	}

	logging.Info(vdfContent["libraryfolders"]["2"]["apps"]["427520"])
	return steamPath
}
