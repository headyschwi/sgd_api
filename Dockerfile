# Use the official MySQL image as the base image
FROM mysql:latest

# Arguments that can be passed at build time
ARG MYSQL_ROOT_PASSWORD
ARG MYSQL_DATABASE

# Set the root password and database for the MySQL server
ENV MYSQL_ROOT_PASSWORD=${MYSQL_ROOT_PASSWORD}
ENV MYSQL_DATABASE=${MYSQL_DATABASE}

# Expose the default MySQL port
EXPOSE 3306
