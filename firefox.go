package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/frioux/leatherman/pkg/mozlz4"
)

const SESSIONBACKUPPATH = "sessionstore-backups/recovery.jsonlz4"

func readMozillaRecoveryBackup(filename string) (*io.Reader, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	r, err := mozlz4.NewReader(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Couldn't create reader: %s\n", err)
		os.Exit(1)
	}
	return &r, nil
}

type MRBEntry struct {
	URL string `json:"url"`
}

type MRBTab struct {
	Entries []MRBEntry `json:"entries"`
}

type MRBWindow struct {
	Tabs []MRBTab `json:"tabs"`
}

type MozillaRecoveryBackup struct {
	Windows []MRBWindow `json:"windows"`
}

func collectUrls(profileDir string) ([]string, error) {
	filename := path.Join(profileDir, SESSIONBACKUPPATH)

	recoveryBackupContents, err := readMozillaRecoveryBackup(filename)
	if err != nil {
		return nil, err
	}
	var mrb MozillaRecoveryBackup
	err = json.NewDecoder(*recoveryBackupContents).Decode(&mrb)
	if err != nil {
		return nil, err
	}
	result := make([]string, 0)
	for w := range mrb.Windows {
		for t := range mrb.Windows[w].Tabs {
			for e := range mrb.Windows[w].Tabs[t].Entries {
				result = append(result, mrb.Windows[w].Tabs[t].Entries[e].URL)
			}
		}
	}
	return result, nil
}
