package main

import (
	"log"

	"github.com/IanMinash/pgbackup/cmd"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("An error occurred loading the .env file; %v\n", err)
	}

	cmd.Execute()
}
