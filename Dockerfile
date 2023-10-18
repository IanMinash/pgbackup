FROM golang:latest

RUN apt update && apt install -y cron postgresql-client postgresql-client-common inetutils-ping

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download 

COPY backup-cron /etc/cron.d/backup-cron

COPY . .

RUN go build -o /root/pgbkp .

RUN chmod 0664 /etc/cron.d/backup-cron && crontab /etc/cron.d/backup-cron

ENTRYPOINT ["cron", "-f"]