#!/bin/bash
set -e

# Variables
APP_NAME="gnotifier"
MAIN_FOLDER="main"

# Define paths
APP_DIR="$HOME/app"

# install dependencies
sudo apt install -y cron make ufw wget
wget https://go.dev/dl/go1.24.1.linux-amd64.tar.gz
sudo rm -rf /usr/local/go && sudo tar -C /usr/local -xzf go1.24.1.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin
go version
sudo rm -f go1.24.1.linux-amd64.tar.gz

# configure firewall
sudo ufw default deny incoming
sudo ufw allow OpenSSH
sudo ufw enable

# install the app
mkdir -p $APP_DIR
cd ../$MAIN_FOLDER
make
sudo cp ../dist/$APP_NAME $APP_DIR
sudo chmod +x $APP_DIR/$APP_NAME
sudo cp config_j/prod.config_j.json $APP_DIR/config_j.json

# Create a cron job to run the app every hour with the -config_j flag
echo "Creating cron job to run $APP_NAME every hour"
echo "run crontab -e"
echo "add to bottom of file : @hourly $APP_DIR/$APP_NAME -config $APP_DIR/config.json"
echo "run sudo systemctl enable cron"
echo "run sudo systemctl restart cron"

echo "Setup complete. $APP_NAME will run every hour."
