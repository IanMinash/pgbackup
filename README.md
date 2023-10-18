# pgbackup

pgbackup is a Go tool that takes a backup of a postgres cluster periodically and saves it into an S3 (Minio) server.


## Configurations
1. Should take a list of clusters that will be backed up
2. Should take in configurations for the S3 server, with permissions to create and delete objects from a bucket.