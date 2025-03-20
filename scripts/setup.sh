#!/bin/bash
set -e

# Variables
SERVER_NAME="gnotifier-server"
CRON_NAME="gnotifier-cron"
APP_DIR="$HOME/app"
LOG_DIR="$HOME/app/logs"
CONFIG_DIR="$HOME/app/config"

# Colors for messages
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to display messages
print_message() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if script is run as root
if [ "$EUID" -eq 0 ]; then 
    print_error "This script should not be run as root"
    exit 1
fi

# Install dependencies
print_message "Installing dependencies..."
sudo apt update
sudo apt install -y cron make ufw wget postgresql postgresql-contrib

# Install Go
print_message "Installing Go..."
wget https://go.dev/dl/go1.24.1.linux-amd64.tar.gz
sudo rm -rf /usr/local/go && sudo tar -C /usr/local -xzf go1.24.1.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin
go version
sudo rm -f go1.24.1.linux-amd64.tar.gz

# Configure firewall
print_message "Configuring firewall..."
sudo ufw default deny incoming
sudo ufw allow OpenSSH
sudo ufw allow 8080/tcp  # For web server
sudo ufw enable

# Create application directories
print_message "Creating application directories..."
mkdir -p $APP_DIR
mkdir -p $LOG_DIR
mkdir -p $CONFIG_DIR

# Build applications
print_message "Building applications..."
cd ..
make build

# Install applications
print_message "Installing applications..."
sudo cp dist/$SERVER_NAME $APP_DIR/
sudo cp dist/$CRON_NAME $APP_DIR/
sudo chmod +x $APP_DIR/$SERVER_NAME
sudo chmod +x $APP_DIR/$CRON_NAME

# Copy configuration files
print_message "Copying configuration files..."
sudo cp server/config.json $CONFIG_DIR/server_config.json
sudo cp cron/config.json $CONFIG_DIR/cron_config.json

# Create systemd service for the server
print_message "Configuring systemd service for the server..."
cat << EOF | sudo tee /etc/systemd/system/gnotifier-server.service
[Unit]
Description=UQAM Grade Notifier Server
After=network.target postgresql.service

[Service]
Type=simple
User=$USER
WorkingDirectory=$APP_DIR
ExecStart=$APP_DIR/$SERVER_NAME -config $CONFIG_DIR/server_config.json
Restart=always
RestartSec=5
StandardOutput=append:/var/log/gnotifier/server.log
StandardError=append:/var/log/gnotifier/server.error.log

[Install]
WantedBy=multi-user.target
EOF

# Create log directory for systemd
sudo mkdir -p /var/log/gnotifier
sudo chown -R $USER:$USER /var/log/gnotifier

# Create cron job for grade checking
print_message "Configuring CRON job..."
(crontab -l 2>/dev/null | grep -v "gnotifier-cron"; echo "0 * * * * $APP_DIR/$CRON_NAME -config $CONFIG_DIR/cron_config.json >> $LOG_DIR/check_grades.log 2>&1") | crontab -

# Enable and start services
print_message "Enabling and starting services..."
sudo systemctl daemon-reload
sudo systemctl enable gnotifier-server
sudo systemctl start gnotifier-server
sudo systemctl enable cron
sudo systemctl restart cron

# Set up PostgreSQL
print_message "Setting up PostgreSQL..."
sudo -u postgres psql -c "CREATE DATABASE uqam_grade_notifier;" || true
sudo -u postgres psql -c "CREATE USER gnotifier WITH PASSWORD 'your_password';" || true
sudo -u postgres psql -c "GRANT ALL PRIVILEGES ON DATABASE uqam_grade_notifier TO gnotifier;" || true

# Set permissions
print_message "Setting permissions..."
sudo chown -R $USER:$USER $APP_DIR
sudo chown -R $USER:$USER $LOG_DIR
sudo chown -R $USER:$USER $CONFIG_DIR

print_message "Installation completed!"
print_warning "Don't forget to:"
echo "1. Modify configuration files in $CONFIG_DIR/"
echo "2. Update database credentials in configuration files"
echo "3. Configure Gmail credentials for notifications"
echo "4. Check logs in $LOG_DIR/ and /var/log/gnotifier/"
echo "5. Restart the server after modifying configuration:"
echo "   sudo systemctl restart gnotifier-server"
