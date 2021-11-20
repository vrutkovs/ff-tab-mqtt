package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"

	"github.com/pierrec/lz4"
)

const SESSIONBACKUPPATH = "sessionstore-backups/recovery.jsonlz4"

func readMozillaRecoveryBackup(filename string) ([]byte, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	size := make([]byte, 4)
	file.Read(size)

	src, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	out := make([]byte, len(src)*10) // XXX should make this use size
	_, err = lz4.UncompressBlock(src, out)
	if err != nil {
		return nil, err
	}
	return out, nil
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
	err = json.Unmarshal(recoveryBackupContents, &mrb)
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
