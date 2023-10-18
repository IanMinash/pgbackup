package main

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/IanMinash/pgbackup/pgbackup"
	"github.com/joho/godotenv"
)

func main() {
	ctx := context.Background()

	err := godotenv.Load()
	if err != nil {
		log.Fatal("An error occurred loading the .env file")
	}
	clustersContent, err := os.ReadFile("./clusters.json")
	if err != nil {
		log.Fatal(err)
	}

	var clusters []pgbackup.ClusterConfig
	err = json.Unmarshal(clustersContent, &clusters)
	if err != nil {
		log.Fatal(err)
	}

	for _, cluster := range clusters {
		archivePath, err := pgbackup.BackupCluster(ctx, &cluster)
		if err != nil {
			log.Fatal(err)
		}

		info, err := pgbackup.UploadArchive(ctx, archivePath)
		if err != nil {
			log.Fatalln(err)
		}
		log.Printf("Version %s of %s has been uploaded!", info.VersionID, info.Key)
	}

	log.Printf("Backups created and uploaded successfully")
}
