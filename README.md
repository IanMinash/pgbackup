# pgbackup

pgbackup is a Go tool that takes a backup of a postgres cluster periodically and saves it into an S3 (Minio) server.

## Prerequisites

This package assumes that the `postgresql-client` package is intalled as it relies on `pg_basebackup` to take a snapshot of a cluster.



## Configurations
1. Should take a list of clusters that will be backed up.
    Create a `clusters.json` file in the current directory with the following format.
    ```json
    [
        {
            "ClusterName": "",
            "Host": "",
            "Port": 5432,
            "Password": "",
            "Username": ""
        },
        ...
    ]
    ```
2. Should take in configurations for the S3 server, with permissions to create objects in a bucket. 
    Create a `.env` file in the current directory with the following.
    ```shell
    # Directory where the backups would be stored
    BACKUP_DIR='./backups'
    S3_ENDPOINT=''
    S3_ACCESS_KEY=''
    S3_SECRET_KEY=''
    S3_BUCKET_NAME=''
    ```

## Running the Docker container
To run the Docker container, use the following command:

```shell
docker run -v /path/to/my/.env:/root/.env -v /path/to/my/clusters.json:/root/clusters.json --network=host ghcr.io/ianminash/pgbackup
```