#!/bin/bash
# Deployment script for AWS EC2

set -e

echo "=== CMS Deployment Script ==="

# Check if Docker is installed
if ! command -v docker &> /dev/null; then
    echo "Installing Docker..."
    sudo yum update -y || sudo apt-get update -y
    sudo yum install -y docker || sudo apt-get install -y docker.io
    sudo systemctl start docker
    sudo systemctl enable docker
    sudo usermod -aG docker $USER
    echo "Docker installed. Please log out and back in, then run this script again."
    exit 0
fi

# Check if Docker Compose is installed
if ! command -v docker-compose &> /dev/null; then
    echo "Installing Docker Compose..."
    sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
    sudo chmod +x /usr/local/bin/docker-compose
fi

# Create .env file if it doesn't exist
if [ ! -f .env ]; then
    echo "Creating .env file..."
    cp .env.example .env

    # Generate a random JWT secret
    JWT_SECRET=$(openssl rand -base64 32)
    sed -i "s/your-super-secret-key-change-this/$JWT_SECRET/" .env

    # Get EC2 public IP
    EC2_IP=$(curl -s http://169.254.169.254/latest/meta-data/public-ipv4 2>/dev/null || echo "localhost")
    sed -i "s/your-ec2-ip/$EC2_IP/" .env

    echo "Generated .env file with JWT_SECRET and EC2 IP: $EC2_IP"
    echo "IMPORTANT: Change ADMIN_PASSWORD in .env before proceeding!"
fi

# Build and start
echo "Building and starting the CMS..."
docker-compose up -d --build

echo ""
echo "=== Deployment Complete ==="
echo "CMS is running on port 8080"
echo ""
echo "Default credentials: admin / admin123"
echo "IMPORTANT: Change the admin password after first login!"
echo ""
echo "To view logs: docker-compose logs -f"
echo "To stop: docker-compose down"
