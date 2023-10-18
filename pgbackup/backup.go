package pgbackup

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

type ClusterConfig struct {
	ClusterName string
	Host        string
	Port        uint16
	Username    string
	Password    string
}

func BackupCluster(config *ClusterConfig) error {
	backupDir := os.Getenv("BACKUP_DIR")
	clusterBackupDir := fmt.Sprintf("%s/%s", backupDir, config.ClusterName)
	clusterBackupArchive := fmt.Sprintf("%s.gzip", clusterBackupDir)

	// Do the backup using pg_basebackup
	backupCmd := exec.Command("pg_basebackup", "-h", config.Host, "-p", fmt.Sprint(config.Port), "-U", config.Username, "-D", clusterBackupDir, "-w")
	backupCmd.Env = append(backupCmd.Env, fmt.Sprintf("PGPASSWORD=%s", config.Password))
	err := backupCmd.Run()
	if err != nil {
		return err
	}

	// Backup was succesfully created under backupDir/clusterName
	// Now we want to zip it into a single archive
	archiveCmd := exec.Command("tar", "-czf", clusterBackupArchive, clusterBackupDir)
	err = archiveCmd.Run()
	if err != nil {
		return err
	}
	log.Printf("Archive for cluster %s has been successfully created.\n", config.ClusterName)

	// Clean  up, by deleting the original backup directory
	delDirCmd := exec.Command("rm", "-r", clusterBackupDir)
	err = delDirCmd.Run()
	if err != nil {
		return err
	}
	return nil
}
