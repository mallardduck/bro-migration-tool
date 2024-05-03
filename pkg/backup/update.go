package backup

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
)

func UpdateLocalIntoBackup(jsonPath string, backupPath string, newBackupName string) {
	jsonData := readLocalClusterJson(jsonPath)

	// TODO: read the backupPath tar.gz file and insert jsonData to update "clusters.management.cattle.io#v3/local.json"
	// But since tar cannot be updated/(to replace data) we just have to make a new one
	// 1. extract current tar into a temp folder,
	tempBackupPath, err := prepareTempBackupDir(backupPath)
	if err != nil {
		log.Fatalf("Error preparing updated backup: %v", err)
	}
	logrus.Info(*tempBackupPath)
	// 2. then update the local cluster data in the temp backup,
	err = updateLocalClusterInBackup(*tempBackupPath, jsonData)
	if err != nil {
		log.Fatalf("Error preparing updated backup: %v", err)
	}
	err = repackBackupFile(*tempBackupPath, newBackupName)
	if err != nil {
		log.Fatalf("Error repacking backup: %v", err)
	}
}

func readLocalClusterJson(jsonPath string) []byte {
	jsonData, err := os.ReadFile(jsonPath)
	if err != nil {
		log.Fatalf("Error reading JSON file: %v", err)
	}

	// Parse the JSON data into an object
	var jsonObject interface{}
	err = json.Unmarshal(jsonData, &jsonObject)
	if err != nil {
		log.Fatalf("Error parsing JSON data: %v", err)
	}
	jsonData, err = json.Marshal(jsonObject)
	return jsonData
}

func prepareTempBackupDir(backupPath string) (*string, error) {
	tmpDir, err := os.MkdirTemp("", "bro-migrate")
	defer os.RemoveAll(tmpDir) // clean up

	if err != nil {
		return nil, fmt.Errorf("error creating temp path: %v", err)
	}
	logrus.Infof("Created temp path at: %v", tmpDir)

	tarReader, err := OpenTarGzReader(backupPath)
	if err != nil {
		return nil, fmt.Errorf("error creating tar reader: %v", err)
	}
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break // End of archive
		}
		if err != nil {
			return nil, fmt.Errorf("error using tar reader: %v", err)
		}
		currentPath := path.Join(tmpDir, header.Name)

		if header.FileInfo().IsDir() {
			logrus.Infof("Attempting to create folder: %v", currentPath)
			err = os.Mkdir(currentPath, 0700)
			if err != nil {
				return nil, fmt.Errorf("cannot MkdirAll: %v", err)
			}
		} else {
			logrus.Infof("Attempting to create file: %v", currentPath)
			newFile, err := os.Create(currentPath)
			if err != nil {
				return nil, fmt.Errorf("cannot create file: %v", err)
			}
			if _, err := io.Copy(newFile, tarReader); err != nil {
				return nil, fmt.Errorf("cannot extract file: %v", err)
			}
		}
	}
	return &tmpDir, nil
}

func updateLocalClusterInBackup(backupPath string, data []byte) error {
	localClusterFilePath := path.Join(backupPath, "clusters.management.cattle.io#v3/local.json")
	localClusterFile, err := os.OpenFile(localClusterFilePath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer localClusterFile.Close()
	_, err = localClusterFile.Write(data)
	if err != nil {
		return err
	}
	logrus.Info("Updated local cluster.")

	return nil
}

func repackBackupFile(backupPath string, name string) error {
	newArchiveName := fmt.Sprintf("%s.tar.gz", name)
	logrus.Infof("Recompressing backup: %v", newArchiveName)
	cwdPath, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	newArchiveFile, err := os.Create(path.Join(cwdPath, newArchiveName))
	defer newArchiveFile.Close()
	// gzip writer
	gw := gzip.NewWriter(newArchiveFile)
	defer gw.Close()
	// tar writer
	tw := tar.NewWriter(gw)
	defer tw.Close()

	walkFunc := func(currPath string, info os.FileInfo, err error) error {
		if currPath == backupPath {
			return nil
		}
		if err != nil {
			return fmt.Errorf("error in walkFunc for %v: %v", currPath, err)
		}
		hdr, err := tar.FileInfoHeader(info, "")
		if err != nil {
			return fmt.Errorf("error creating header for %v: %v", info.Name(), err)
		}
		// backupPath could be /var/tmp/folders/backup/authconfigs.management.cattle.io/adfs.json
		// we need to include only authconfigs.management.cattle.io onwards, so get relative path
		relativePath, err := filepath.Rel(backupPath, currPath)
		if err != nil {
			return fmt.Errorf("error getting relative path for %v: %v", info.Name(), err)
		}
		hdr.Name = filepath.Join(relativePath)
		if err := tw.WriteHeader(hdr); err != nil {
			return fmt.Errorf("error writing header for %v: %v", info.Name(), err)
		}
		if info.IsDir() {
			return nil
		}
		fInfo, err := os.Open(currPath)
		if err != nil {
			return fmt.Errorf("error opening %v: %v", info.Name(), err)
		}
		if _, err := io.Copy(tw, fInfo); err != nil {
			return fmt.Errorf("error copying %v: %v", info.Name(), err)
		}
		return fInfo.Close()
	}
	filepath.Walk(backupPath, walkFunc)

	return nil
}
