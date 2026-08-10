#!/usr/bin/env bash

set -e

# Ensure both variables are provided
if [ -n "$URL" ] && [ -n "$GCS_BUCKET" ]; then
    # Perform the database backup
    pg_dump -v "$URL" > /usr/src/app/backup.sql

    # Upload directly to GCS (Workload Identity handles authentication automatically)
    TIMESTAMP=$(date +"%Y%m%d_%H%M%S")
    gsutil cp /usr/src/app/backup.sql "gs://${GCS_BUCKET}/backup_${TIMESTAMP}.sql"

    echo "Backup successfully uploaded to gs://${GCS_BUCKET}/backup_${TIMESTAMP}.sql"
else
    echo "Error: URL and GCS_BUCKET environment variables must be set." >&2
    exit 1
fi
