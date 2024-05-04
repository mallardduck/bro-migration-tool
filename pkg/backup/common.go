package backup

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"os"
)

func OpenTarGzReader(filePath string) (*tar.Reader, error) {
	gzipFile, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("error opening backup tar gzip file: %v", err)
	}
	// setup stream readers
	gzipReader, err := gzip.NewReader(gzipFile)
	if err != nil {
		return nil, fmt.Errorf("error creating gzip reader: %v", err)
	}
	defer gzipReader.Close()
	tarReader := tar.NewReader(gzipReader)

	return tarReader, nil
}

func FetchLocalClusterFromBackup(backupPath string, jsonPath string) (map[string]interface{}, error) {
	jsonData := make(map[string]interface{})
	localClusterBytes, err := extractLocalClusterFromBackup(backupPath)
	if err != nil {
		return jsonData, err
	}

	err = json.Unmarshal(localClusterBytes, &jsonData)
	if err != nil {
		return jsonData, err
	}

	return jsonData, nil
}
