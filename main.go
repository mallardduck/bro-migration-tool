package main

import (
	"errors"
	"fmt"
	"github.com/mallardduck/bro-migration-tool/pkg/migrate"
	"os"

	"github.com/mallardduck/bro-migration-tool/pkg/backup"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var (
	// Version represents the current version of the chart build scripts
	Version = "v0.0.0-dev"
	// GitCommit represents the latest commit when building this script
	GitCommit = "HEAD"

	BackupFilePath       string
	NewBackupFilename    string
	LocalClusterFilePath = "./local.json"
)

func main() {
	if len(os.Getenv("DEBUG")) > 0 {
		logrus.SetLevel(logrus.DebugLevel)
	}
	app := cli.NewApp()
	app.Name = "bro-migration-tool"
	app.Version = fmt.Sprintf("%s (%s)", Version, GitCommit)

	fileFlag := cli.StringFlag{
		Name:        "file,f",
		Usage:       "The file path",
		Required:    true,
		Destination: &BackupFilePath,
	}

	fileOutFlag := cli.StringFlag{
		Name:        "out,o",
		Usage:       "The file out path",
		Required:    true,
		Destination: &NewBackupFilename,
	}

	app.Commands = []cli.Command{
		{
			Name:   "pull-local",
			Usage:  "Pull the local cluster out of the backup file",
			Action: pullLocalJson,
			Flags:  []cli.Flag{fileFlag},
		},
		{
			Name:   "push-local",
			Usage:  "Push the local cluster out of the backup file",
			Action: pushLocalJson,
			Flags:  []cli.Flag{fileFlag, fileOutFlag},
		},
		{
			Name:   "k3s-rke2",
			Usage:  "Prepare a K3S origin backup for RKE2 restore.",
			Action: transformK3sBackupForRke2,
			Flags:  []cli.Flag{fileFlag, fileOutFlag},
		},
		{
			Name:   "rke2-k3s",
			Usage:  "Prepare a RKE2 origin backup for k3s restore.",
			Action: transformRke2BackupForK3s,
			Flags:  []cli.Flag{fileFlag, fileOutFlag},
		},
	}

	if err := app.Run(os.Args); err != nil {
		logrus.Fatal(err)
	}
}

func pullLocalJson(c *cli.Context) {
	// Verify the file exists at the path...
	if _, err := os.Stat(BackupFilePath); errors.Is(err, os.ErrNotExist) {
		// path/to/whatever does not exist
		logrus.Fatal("The backup archive file doesn't exist.")
	}
	// TODO: actually do something with these results
	_, _ = backup.PullLocalFromBackup(BackupFilePath, LocalClusterFilePath)
}

func pushLocalJson(c *cli.Context) {
	// Verify the BackupFilePath exists...
	if _, err := os.Stat(BackupFilePath); errors.Is(err, os.ErrNotExist) {
		// path/to/whatever does not exist
		logrus.Fatal("The backup archive file doesn't exist.")
	}
	// Verify the LocalClusterFilePath exists...
	if _, err := os.Stat(LocalClusterFilePath); errors.Is(err, os.ErrNotExist) {
		// path/to/whatever does not exist
		logrus.Fatal("The local cluster file doesn't exist.")
	}
	localClusterData := backup.ReadLocalClusterJson(LocalClusterFilePath)
	backup.UpdateLocalIntoBackup(localClusterData, BackupFilePath, NewBackupFilename)
	logrus.Infoln(NewBackupFilename)
}

func transformK3sBackupForRke2(c *cli.Context) {
	// Verify the file exists at the path...
	if _, err := os.Stat(BackupFilePath); errors.Is(err, os.ErrNotExist) {
		// path/to/whatever does not exist
		logrus.Fatal("The backup archive file doesn't exist.")
	}
	localClusterData, err := backup.FetchLocalClusterFromBackup(BackupFilePath, LocalClusterFilePath)
	if err != nil {
		logrus.Fatal("Failed to fetch the local cluster from backup file")
	}
	newClusterData := migrate.K3sRancherToRke2Rancher(localClusterData)
	backup.UpdateLocalIntoBackup(newClusterData, BackupFilePath, NewBackupFilename)
}

func transformRke2BackupForK3s(c *cli.Context) {
	// Verify the file exists at the path...
	if _, err := os.Stat(BackupFilePath); errors.Is(err, os.ErrNotExist) {
		// path/to/whatever does not exist
		logrus.Fatal("The backup archive file doesn't exist.")
	}
	localClusterData, err := backup.FetchLocalClusterFromBackup(BackupFilePath, LocalClusterFilePath)
	if err != nil {
		logrus.Fatal("Failed to fetch the local cluster from backup file")
	}
	// TODO: do the thing to migrate
	newClusterData := migrate.Rke2RancherToK3sRancher(localClusterData)
	backup.UpdateLocalIntoBackup(newClusterData, BackupFilePath, NewBackupFilename)
}
