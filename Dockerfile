# Use the official MySQL image as the base image
FROM mysql:latest

# Arguments that can be passed at build time
ARG MYSQL_ROOT_PASSWORD
ARG MYSQL_DATABASE
ARG MYSQL_PORT
ARG MYSQL_CHARSET=utf8mb4
ARG MYSQL_COLLATION=utf8mb4_unicode_ci

# Set the root password for the MySQL server
ENV MYSQL_ROOT_PASSWORD=${MYSQL_ROOT_PASSWORD}

# Create a new database
ENV MYSQL_DATABASE=${MYSQL_DATABASE}

# Set default port for MySQL
ENV MYSQL_PORT=${MYSQL_PORT}

# Set the character set and collation for the database
ENV MYSQL_CHARSET=${MYSQL_CHARSET}
ENV MYSQL_COLLATION=${MYSQL_COLLATION}

# Copy any SQL scripts to initialize the database
COPY init.sql /docker-entrypoint-initdb.d/

# Expose the MySQL port
EXPOSE ${MYSQL_PORT}
