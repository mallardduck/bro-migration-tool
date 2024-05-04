package backup

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
)

func PullLocalFromBackup(backupPath string, jsonPath string) (bool, error) {
	localClusterBytes, err := extractLocalClusterFromBackup(backupPath)
	if err != nil {
		return false, err
	}

	var jsonData bytes.Buffer
	json.Indent(&jsonData, localClusterBytes, "", "\t")

	err = saveLocalClusterJson(jsonPath, jsonData)
	if err != nil {
		return false, err
	}

	return true, nil
}

func extractLocalClusterFromBackup(backupPath string) ([]byte, error) {
	tarReader, err := OpenTarGzReader(backupPath)
	if err != nil {
		return nil, fmt.Errorf("error creating tar reader: %v", err)
	}

	var localClusterBytes []byte
	// loop through the files in the tar archive and find the desired file
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break // end of archive
		} else if err != nil {
			return nil, fmt.Errorf("error reading tar header: %v", err)
		}

		// check if the file path matches the desired path
		if header.Name == "clusters.management.cattle.io#v3/local.json" {
			// extract the contents of the file
			content, err := io.ReadAll(tarReader)
			if err != nil {
				return nil, fmt.Errorf("error reading tar content: %v", err)
			}

			// do something with the contents of the file
			localClusterBytes = content
		}
	}
	return localClusterBytes, nil
}

func saveLocalClusterJson(jsonPath string, jsonData bytes.Buffer) error {
	outFile, err := os.Create(jsonPath)
	if err != nil {
		panic(err)
	}
	defer outFile.Close()
	_, err = outFile.Write(jsonData.Bytes())
	if err != nil {
		return err
	}
	return nil
}
