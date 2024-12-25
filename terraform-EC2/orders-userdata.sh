#!/bin/bash

# Update and install essential tools
echo "Updating system and installing prerequisites..."
apt-get update -y && apt-get upgrade -y
apt-get install -y apt-transport-https ca-certificates curl software-properties-common gnupg lsb-release unzip

# Install Docker
echo "Installing Docker..."
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | gpg --dearmor -o /usr/share/keyrings/docker-archive-keyring.gpg
echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable" | tee /etc/apt/sources.list.d/docker.list > /dev/null
apt-get update -y
apt-get install -y docker-ce docker-ce-cli containerd.io
systemctl start docker
systemctl enable docker

# Add user to Docker group
echo "Adding user to Docker group..."
usermod -aG docker ubuntu

# Clean up unused packages
apt-get autoremove -y
apt-get clean

# Verify Docker installation
echo "Docker version:"
docker --version || { echo "Docker installation failed!"; exit 1; }

# Clone repository
REPO_URL="https://github.com/sujalsharmaa/go-grpc-graphql-microservices.git"
REPO_DIR="go-grpc-graphql-microservices"

echo "Cloning repository..."
git clone -b features/cloud-and-devops-features "$REPO_URL" || { echo "Failed to clone repository!"; exit 1; }


# Build Docker image
echo "Building Docker image..."
docker build -f order/app.dockerfile -t orders . || { echo "Docker build failed!"; exit 1; }

# Run Docker container
echo "Running Docker container..."

echo "Setup complete! You can access the application on port 80."


docker run -d -p 80:8080
      -e DATABASE_URL="postgres.orders.backend.in" \
      -e ACCOUNT_SERVICE_URL="http://backend.accounts.com" \
      -e CATALOG_SERVICE_URL="http://backend.catalog.com"
      orders || { echo 'Failed to start Docker container!'; exit 1; }