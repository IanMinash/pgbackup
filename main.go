package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/IanMinash/pgbackup/pgbackup"
	"github.com/joho/godotenv"
)

func main() {
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

		err = pgbackup.BackupCluster(&cluster)
		if err != nil {
			log.Fatal(err)
		}
	}

	backupDir := os.Getenv("BACKUP_DIR")

	fmt.Printf("Backups created successfully, check %s\n", backupDir)
}
