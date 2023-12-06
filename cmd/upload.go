package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/IanMinash/pgbackup/pgbackup"
	"github.com/spf13/cobra"
)

var directory string
var archiveFile string
var retainDays uint16

func init() {
	uploadCmd.Flags().StringVarP(&directory, "directory", "d", "./backups", "Directory containing cluster archives to upload")
	uploadCmd.Flags().StringVarP(&archiveFile, "archive-file", "f", "", "A cluster archive to upload")
	uploadCmd.Flags().Uint16VarP(&retainDays, "retain-days", "p", 7, "The number of days this archive should be retained for.")
	rootCmd.AddCommand(uploadCmd)
}

var uploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "Upload a zipped file or files in a directory to S3",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()
		if archiveFile != "" {
			info, err := pgbackup.UploadArchive(ctx, archiveFile, retainDays)
			if err != nil {
				log.Printf("Error occurred while uploading archive %s: %v\n", archiveFile, err)
				return
			}
			log.Printf("Version %s of %s has been uploaded!", info.VersionID, info.Key)
		} else {
			fileEntries, err := os.ReadDir(directory)
			if err != nil {
				log.Printf("An error occurred while reading %s: %v\n", directory, err)
				return
			}
			var wg sync.WaitGroup

			for _, entry := range fileEntries {
				if entry.IsDir() {
					continue
				}
				fileName := fmt.Sprintf("%s/%s", directory, entry.Name())
				if strings.HasSuffix(fileName, ".zip") {
					wg.Add(1)

					go func(archivePath string) {
						defer wg.Done()

						info, err := pgbackup.UploadArchive(ctx, archivePath, retainDays)
						if err != nil {
							log.Printf("Error occurred while uploading archive %s: %v\n", archivePath, err)
							return
						}
						log.Printf("Version %s of %s has been uploaded!", info.VersionID, info.Key)
					}(fileName)
				}
				entry.Name()
			}
			wg.Wait()
		}
	},
}
