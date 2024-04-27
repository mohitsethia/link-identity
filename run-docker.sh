#!/bin/bash

docker-compose up -d

# Build the Docker image and capture the output
output=$(docker build --build-arg APP_NAME=link-identity-api .)

# Use grep to find the "Successfully built" line and awk to extract the image ID
image_id=$(echo "$output" | grep "Successfully built" | awk '{print $3}')

# Run a container using the captured image ID
docker run -p 8000:8000 $image_id