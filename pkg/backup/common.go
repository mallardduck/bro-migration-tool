package backup

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
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

func FetchLocalClusterFromBackup(backupPath string, jsonPath string) (unstructured.Unstructured, error) {
	clusterObject := unstructured.Unstructured{}
	localClusterBytes, err := extractLocalClusterFromBackup(backupPath)
	if err != nil {
		return clusterObject, err
	}

	err = clusterObject.UnmarshalJSON(localClusterBytes)
	return clusterObject, err
}
