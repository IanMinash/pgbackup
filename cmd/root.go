package cmd

import (
	"encoding/json"
	"log"
	"os"
	"sync"

	"github.com/IanMinash/pgbackup/pgbackup"
	"github.com/spf13/cobra"
)

var ClustersInputUsage string = `Json file containing the clusters. It should be a JSON array with each element having the following keys: 
ClusterName (string), Host (string), Port (int), Password (string), Username (string)
`

var rootCmd = &cobra.Command{
	Use:   "pgbackup",
	Short: "pgbackup is a Go tool that takes a backup of a postgres cluster periodically and saves it into an S3 (Minio) server.",
	Run: func(cmd *cobra.Command, args []string) {
		clustersContent, err := os.ReadFile("./clusters.json")
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

				info, err := pgbackup.UploadArchive(ctx, archivePath)
				if err != nil {
					log.Printf("Error occurred while uploading backup for cluster %s: %v\n", cl.ClusterName, err)
					return
				}
				log.Printf("Version %s of %s has been uploaded!", info.VersionID, info.Key)
			}(cluster)
		}

		wg.Wait()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
