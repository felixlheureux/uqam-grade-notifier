#!/bin/bash
set -e

# Variables
APP_NAME="gnotifier"

# Define paths
APP_DIR="/usr/local/bin"
CONFIG_DIR="/etc/$APP_NAME"
LOG_DIR="/var/log/$APP_NAME"
CRON_FILE="/etc/cron.d/$APP_NAME"

# Create necessary directories
sudo mkdir -p $CONFIG_DIR $LOG_DIR

# Build the Go application
make build

# Move executable to the appropriate location
sudo mv $APP_NAME $APP_DIR/$APP_NAME
sudo chmod +x $APP_DIR/$APP_NAME

# Copy config and data files
sudo cp config.json $CONFIG_DIR/config.json
sudo cp grades.json $CONFIG_DIR/grades.json

# Create a cron job to run the app every hour
echo "0 * * * * root $APP_DIR/$APP_NAME >> $LOG_DIR/$APP_NAME.log 2>&1" | sudo tee $CRON_FILE
sudo chmod 644 $CRON_FILE

# Ensure cron service is running
sudo systemctl restart cron

echo "Setup complete. $APP_NAME will run every hour."
