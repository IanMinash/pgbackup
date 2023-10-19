package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"sync"

	"github.com/IanMinash/pgbackup/pgbackup"
	"github.com/joho/godotenv"
)

func main() {
	ctx := context.Background()

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("An error occurred loading the .env file; %v\n", err)
	}

	clustersContent, err := os.ReadFile("./clusters.json")
	if err != nil {
		log.Fatalln(err)
	}

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

	log.Printf("Backups created and uploaded successfully")
}
