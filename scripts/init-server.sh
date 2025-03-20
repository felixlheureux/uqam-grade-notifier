#!/bin/bash
set -e

# Check if Podman is installed
if ! command -v podman &> /dev/null; then
    echo "Installing Podman..."
    sudo apt update
    sudo apt install -y podman podman-compose
fi

# Create application directory
sudo mkdir -p /opt/gnotifier
sudo chown -R $USER:$USER /opt/gnotifier

# Create directories for Nginx and backups
sudo mkdir -p /opt/gnotifier/nginx/conf.d
sudo mkdir -p /opt/gnotifier/nginx/ssl
sudo mkdir -p /opt/gnotifier/backups

# Create .env file
cat > /opt/gnotifier/.env << EOF
DB_PASSWORD=${DB_PASSWORD}
GMAIL_APP_PASSWORD=${GMAIL_APP_PASSWORD}
JWT_SECRET=${JWT_SECRET}
ENCRYPTION_KEY=${ENCRYPTION_KEY}
EOF

# Copy configuration files
cp podman-compose.yml /opt/gnotifier/
cp nginx/conf.d/default.conf /opt/gnotifier/nginx/conf.d/

# Start services
cd /opt/gnotifier
podman-compose up -d

# Wait for services to be ready
sleep 10

# Get SSL certificate
podman-compose run --rm certbot

# Restart Nginx to apply SSL configuration
podman-compose restart nginx

# Create systemd service for automatic restart
sudo tee /etc/systemd/system/gnotifier.service << EOF
[Unit]
Description=GNotifier Podman Compose
Requires=podman.service
After=podman.service

[Service]
Type=oneshot
RemainAfterExit=yes
WorkingDirectory=/opt/gnotifier
ExecStart=/usr/bin/podman-compose up -d
ExecStop=/usr/bin/podman-compose down
TimeoutStartSec=0
User=$USER
Group=$USER

[Install]
WantedBy=multi-user.target
EOF

# Enable systemd service
sudo systemctl daemon-reload
sudo systemctl enable gnotifier.service
sudo systemctl start gnotifier.service

# Create systemd service for backups
sudo tee /etc/systemd/system/gnotifier-backup.service << EOF
[Unit]
Description=GNotifier Database Backup
Requires=gnotifier.service
After=gnotifier.service

[Service]
Type=oneshot
WorkingDirectory=/opt/gnotifier
ExecStart=/usr/bin/podman-compose exec -T db pg_dump -U gnotifier uqam_grade_notifier > /opt/gnotifier/backups/backup_$(date +%%Y%%m%%d_%%H%%M%%S).sql
ExecStartPost=/usr/bin/find /opt/gnotifier/backups -type f -mtime +7 -delete
User=$USER
Group=$USER

[Install]
WantedBy=multi-user.target
EOF

# Create timer for daily backups
sudo tee /etc/systemd/system/gnotifier-backup.timer << EOF
[Unit]
Description=Daily GNotifier Database Backup

[Timer]
OnCalendar=daily
Persistent=true

[Install]
WantedBy=timers.target
EOF

# Enable backup timer
sudo systemctl daemon-reload
sudo systemctl enable gnotifier-backup.timer
sudo systemctl start gnotifier-backup.timer

# Configure sub-uid and sub-gid for Podman
echo "$USER:100000:65536" | sudo tee /etc/subuid
echo "$USER:100000:65536" | sudo tee /etc/subgid

# Create Podman configuration directory
mkdir -p ~/.config/containers
cat > ~/.config/containers/containers.conf << EOF
[engine]
runtime = "crun"
events_logger = "file"

[engine.runtimes]
crun = ["/usr/bin/crun", ""]
EOF 