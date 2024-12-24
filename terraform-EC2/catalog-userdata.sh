#!/bin/bash

# Update and install essential tools
echo "Updating system and installing prerequisites..."
sudo apt-get update -y && sudo apt-get upgrade -y
sudo apt-get install -y apt-transport-https ca-certificates curl software-properties-common gnupg lsb-release unzip

# Install Docker
echo "Installing Docker..."
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /usr/share/keyrings/docker-archive-keyring.gpg
echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
sudo apt-get update -y
sudo apt-get install -y docker-ce docker-ce-cli containerd.io
sudo systemctl start docker
sudo systemctl enable docker

# Add user to Docker group
echo "Adding user to Docker group..."
sudo usermod -aG docker ubuntu

# Clean up unused packages
sudo apt-get autoremove -y
sudo apt-get clean

# Verify Docker installation
echo "Docker version:"
docker --version || { echo "Docker installation failed!"; exit 1; }

# Clone repository
REPO_URL="https://github.com/sujalsharmaa/go-grpc-graphql-microservices.git"
REPO_DIR="go-grpc-graphql-microservices"

echo "Cloning repository..."
git clone "$REPO_URL" || { echo "Failed to clone repository!"; exit 1; }

# Navigate to the catalog directory
cd "$REPO_DIR/catalog" || { echo "Failed to navigate to catalog!"; exit 1; }

# Build Docker image
echo "Building Docker image..."
docker build -f app.dockerfile -t catalog . || { echo "Docker build failed!"; exit 1; }

# Run Docker container
echo "Running Docker container..."

echo "Setup complete! You can access the application on port 80."

docker run -d -p 80:8080 -e ELASTICSEARCH_DB_URL=http://backend.elasticSearch.com catalog || { echo 'Failed to start Docker container!'; exit 1; }