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
REPO_URL="https://github.com/AkhilSharma90/go-grpc-graphql-microservices.git"
REPO_DIR="go-grpc-graphql-microservices"

echo "Cloning repository..."
git clone "$REPO_URL" || { echo "Failed to clone repository!"; exit 1; }

# Navigate to the backend-nodejs directory
cd "$REPO_DIR/graphql" || { echo "Failed to navigate to account!"; exit 1; }

# Build Docker image
echo "Building Docker image..."
docker build -f app.dockerfile -t graphql . || { echo "Docker build failed!"; exit 1; }

# Run Docker container
echo "Running Docker container..."

echo "Setup complete! You can access the application on port 80."
 sudo docker run -d -p 80:8080 \
        -e CATALOG_SERVICE_URL="http://backend.catalog.com" \
        -e ACCOUNT_SERVICE_URL="http://backend.accounts.com" \
        -e ORDER_SERVICE_URL="http://backend.orders.com" 
        graphql || { echo 'Failed to start Docker container!'; exit 1; }


