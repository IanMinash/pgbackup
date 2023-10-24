package cmd

import (
	"encoding/json"
	"log"
	"os"
	"sync"

	"github.com/IanMinash/pgbackup/pgbackup"
	"github.com/spf13/cobra"
)

var clustersFile string

func init() {
	backupCmd.Flags().StringVarP(&clustersFile, "input-file", "i", "./clusters.json", ClustersInputUsage)
	rootCmd.AddCommand(backupCmd)
}

var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Create gzip compressed backups of a Postgres cluster",
	Run: func(cmd *cobra.Command, args []string) {

		clustersContent, err := os.ReadFile(clustersFile)
		if err != nil {
			log.Fatalln(err)
		}

		ctx := cmd.Context()

		var clusters []pgbackup.ClusterConfig
		err = json.Unmarshal(clustersContent, &clusters)
		if err != nil {
			log.Fatalln(err)
		}

		var wg sync.WaitGroup

		for _, cluster := range clusters {
			wg.Add(1)

			go func(cl pgbackup.ClusterConfig) {
				defer wg.Done()

				archivePath, err := pgbackup.BackupCluster(ctx, &cl)
				if err != nil {
					log.Printf("Error occurred while generating backup for cluster %s: %v\n", cl.ClusterName, err)
					return
				}

				log.Printf("Backup for cluster %s has been successfully generated: %s", cl.ClusterName, archivePath)
			}(cluster)
		}

		wg.Wait()
	},
}
