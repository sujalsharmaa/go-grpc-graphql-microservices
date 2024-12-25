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
sudo usermod -aG docker "$USER"

# Clean up unused packages
sudo apt-get autoremove -y
sudo apt-get clean

# Verify Docker installation
echo "Docker version:"
docker --version || { echo "Docker installation failed!"; exit 1; }

apt install postgresql-client-16 -y



# Clone repository
REPO_URL="https://github.com/sujalsharmaa/go-grpc-graphql-microservices.git"
REPO_DIR="go-grpc-graphql-microservices"

echo "Cloning repository..."
git clone -b features/cloud-and-devops-features "$REPO_URL" || { echo "Failed to clone repository!"; exit 1; }

# Navigate to the account service directory

# Build Docker image
echo "Building Docker image..."
docker build -f account/app.dockerfile -t accounts . || { echo "Docker build failed!"; exit 1; }

# # Run Docker container
# echo "Running Docker container..."

docker run -d -p 80:8080 -e DATABASE_URL="postgres.accounts.backend.com" accounts || { echo "Failed to start Docker container!"; exit 1; }

# echo "Setup complete! You can access the application on port 80."
