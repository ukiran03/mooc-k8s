FROM gcr.io/google.com/cloudsdktool/google-cloud-cli:alpine

# Install PostgreSQL client (which includes pg_dump)
RUN apk add --no-cache postgresql-client bash

# Set working directory
WORKDIR /usr/src/app

# Copy your backup script into the container
COPY backup.sh /usr/src/app/backup.sh
RUN chmod +x /usr/src/app/backup.sh

# Set the default entrypoint to run the script
ENTRYPOINT ["/usr/src/app/backup.sh"]
