#!/bin/bash
set -e

# Variables
SERVER_NAME="gnotifier-server"
CRON_NAME="gnotifier-cron"
APP_DIR="$HOME/app"
LOG_DIR="$HOME/app/logs"
CONFIG_DIR="$HOME/app/config"

# Couleurs pour les messages
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Fonction pour afficher les messages
print_message() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Vérifier si le script est exécuté en tant que root
if [ "$EUID" -eq 0 ]; then 
    print_error "Ce script ne doit pas être exécuté en tant que root"
    exit 1
fi

# Install dependencies
print_message "Installation des dépendances..."
sudo apt update
sudo apt install -y cron make ufw wget postgresql postgresql-contrib

# Install Go
print_message "Installation de Go..."
wget https://go.dev/dl/go1.24.1.linux-amd64.tar.gz
sudo rm -rf /usr/local/go && sudo tar -C /usr/local -xzf go1.24.1.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin
go version
sudo rm -f go1.24.1.linux-amd64.tar.gz

# Configure firewall
print_message "Configuration du pare-feu..."
sudo ufw default deny incoming
sudo ufw allow OpenSSH
sudo ufw allow 8080/tcp  # Pour le serveur web
sudo ufw enable

# Create application directories
print_message "Création des répertoires de l'application..."
mkdir -p $APP_DIR
mkdir -p $LOG_DIR
mkdir -p $CONFIG_DIR

# Build applications
print_message "Compilation des applications..."
cd ..
make build

# Install applications
print_message "Installation des applications..."
sudo cp dist/$SERVER_NAME $APP_DIR/
sudo cp dist/$CRON_NAME $APP_DIR/
sudo chmod +x $APP_DIR/$SERVER_NAME
sudo chmod +x $APP_DIR/$CRON_NAME

# Copy configuration files
print_message "Copie des fichiers de configuration..."
sudo cp server/config.json $CONFIG_DIR/server_config.json
sudo cp cron/config.json $CONFIG_DIR/cron_config.json

# Create systemd service for the server
print_message "Configuration du service systemd pour le serveur..."
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
print_message "Configuration de la tâche CRON..."
(crontab -l 2>/dev/null | grep -v "gnotifier-cron"; echo "0 * * * * $APP_DIR/$CRON_NAME -config $CONFIG_DIR/cron_config.json >> $LOG_DIR/check_grades.log 2>&1") | crontab -

# Enable and start services
print_message "Activation et démarrage des services..."
sudo systemctl daemon-reload
sudo systemctl enable gnotifier-server
sudo systemctl start gnotifier-server
sudo systemctl enable cron
sudo systemctl restart cron

# Set up PostgreSQL
print_message "Configuration de PostgreSQL..."
sudo -u postgres psql -c "CREATE DATABASE uqam_grade_notifier;" || true
sudo -u postgres psql -c "CREATE USER gnotifier WITH PASSWORD 'your_password';" || true
sudo -u postgres psql -c "GRANT ALL PRIVILEGES ON DATABASE uqam_grade_notifier TO gnotifier;" || true

# Set permissions
print_message "Configuration des permissions..."
sudo chown -R $USER:$USER $APP_DIR
sudo chown -R $USER:$USER $LOG_DIR
sudo chown -R $USER:$USER $CONFIG_DIR

print_message "Installation terminée !"
print_warning "N'oubliez pas de :"
echo "1. Modifier les fichiers de configuration dans $CONFIG_DIR/"
echo "2. Mettre à jour les identifiants de la base de données dans les fichiers de configuration"
echo "3. Configurer les identifiants Gmail pour les notifications"
echo "4. Vérifier les logs dans $LOG_DIR/ et /var/log/gnotifier/"
echo "5. Redémarrer le serveur après avoir modifié la configuration :"
echo "   sudo systemctl restart gnotifier-server"
