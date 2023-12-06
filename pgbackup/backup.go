package pgbackup

import (
	"context"
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
	RetainDays  uint16
}

// BackupCluster runs pg_basebackup on a specified cluster and
// returns the path of the compressed archive.
func BackupCluster(ctx context.Context, config *ClusterConfig) (string, error) {
	backupDir := os.Getenv("BACKUP_DIR")
	clusterBackupDir := fmt.Sprintf("%s/%s", backupDir, config.ClusterName)
	clusterBackupArchive := fmt.Sprintf("%s.gz", clusterBackupDir)

	// Do the backup using pg_basebackup
	backupCmd := exec.CommandContext(ctx, "pg_basebackup", "-z", "-F", "t", "-h", config.Host, "-p", fmt.Sprint(config.Port), "-U", config.Username, "-D", clusterBackupDir, "-w")
	backupCmd.Env = append(backupCmd.Env, fmt.Sprintf("PGPASSWORD=%s", config.Password))
	err := backupCmd.Run()
	if err != nil {
		return "", err
	}

	// Compress into single archive
	archiveCmd := exec.CommandContext(ctx, "tar", "-cf", clusterBackupArchive, "-C", clusterBackupDir, ".")
	err = archiveCmd.Run()
	if err != nil {
		return "", err
	}

	log.Printf("Archive for cluster %s has been successfully created.\n", config.ClusterName)

	// Clean  up, by deleting the original backup directory
	delDirCmd := exec.CommandContext(ctx, "rm", "-r", clusterBackupDir)
	err = delDirCmd.Run()
	if err != nil {
		return "", err
	}

	return clusterBackupArchive, nil
}
