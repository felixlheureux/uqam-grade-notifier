terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 5.0"
    }
  }
}

provider "google" {
  project = var.project_id
  region  = var.region
  zone    = var.zone
}

# Création de l'instance e2-micro
resource "google_compute_instance" "uqam_grade_notifier" {
  name         = "uqam-grade-notifier"
  machine_type = "e2-micro"

  boot_disk {
    initialize_params {
      image = "debian-cloud/debian-12"
    }
  }

  network_interface {
    network = "default"
    access_config {
      // Ephemeral public IP
    }
  }

  metadata = {
    ssh-keys = "${var.ssh_user}:${file(var.ssh_pub_key_path)}"
  }

  tags = ["uqam-grade-notifier"]

  # Script de démarrage pour installer Docker et l'application
  metadata_startup_script = <<-EOF
    #!/bin/bash
    # Installation de Docker
    apt-get update
    apt-get install -y apt-transport-https ca-certificates curl software-properties-common
    curl -fsSL https://download.docker.com/linux/debian/gpg | apt-key add -
    add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/debian $(lsb_release -cs) stable"
    apt-get update
    apt-get install -y docker-ce docker-ce-cli containerd.io

    # Création du répertoire de l'application
    mkdir -p /opt/uqam-grade-notifier
    cd /opt/uqam-grade-notifier

    # Cloner le repository (à remplacer par votre repository)
    git clone https://github.com/felixlheureux/uqam-grade-notifier.git .

    # Construire et démarrer l'application avec Docker
    docker-compose up -d --build
  EOF
}

# Règle de pare-feu pour autoriser le trafic HTTP
resource "google_compute_firewall" "allow_http" {
  name    = "allow-http"
  network = "default"

  allow {
    protocol = "tcp"
    ports    = ["80", "443"]
  }

  source_ranges = ["0.0.0.0/0"]
  target_tags   = ["uqam-grade-notifier"]
} 