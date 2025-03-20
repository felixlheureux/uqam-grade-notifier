# Create a VPC network
resource "google_compute_network" "vpc" {
  name                    = "${var.environment}-${var.project}-vpc"
  auto_create_subnetworks = false
  routing_mode            = "REGIONAL"
  mtu                     = 1460 
}

# Create a public subnet
resource "google_compute_subnetwork" "public_subnet" {
  name          = "${var.environment}-${var.project}-public-subnet"
  region        = var.region
  network       = google_compute_network.vpc.self_link
  ip_cidr_range = "10.0.1.0/24"  
  depends_on    = [google_compute_network.vpc]
}

# Basic Network Firewall Rules
# Allow http
resource "google_compute_firewall" "allow-http" {
  name    = "${var.environment}-${var.project}-allow-http"
  network = "${google_compute_network.vpc.id}"
  allow {
    protocol = "tcp"
    ports    = ["80"]
  }
  direction     = "INGRESS"
  source_ranges = ["0.0.0.0/0"]
  target_tags = ["http"] 
}

# allow https
resource "google_compute_firewall" "allow-https" {
  name    = "${var.environment}-${var.project}-allow-https"
  network = "${google_compute_network.vpc.id}"
  allow {
    protocol = "tcp"
    ports    = ["443"]
  }
  
  direction     = "INGRESS"
  source_ranges = ["0.0.0.0/0"]
  target_tags = ["https"] 
}

# allow ssh
resource "google_compute_firewall" "allow-ssh" {
  name    = "${var.environment}-${var.project}-allow-ssh"
  network = "${google_compute_network.vpc.id}"
  allow {
    protocol = "tcp"
    ports    = ["22"]
  }
  direction     = "INGRESS"
  source_ranges = ["0.0.0.0/0"]
  target_tags = ["ssh"]
}

# Selects the OS for the GCP VM.
data "google_compute_image" "image" {
  family  = "debian-12"
  project = "debian-cloud"
}

# Create a static IP address
resource "google_compute_address" "static-ip-address" {
  name = "${var.environment}-${var.project}-static-ip-address"
}

# Create a single Compute Engine instance
resource "google_compute_instance" "default" {
  name         = "${var.environment}-${var.project}-vm"
  machine_type = var.machine_type
  zone         = "${var.region}-b"
  tags         = ["${var.environment}", "http", "https", "ssh"]

  boot_disk {
    initialize_params {
        image = data.google_compute_image.image.self_link
    }
  }

  network_interface {
    subnetwork = google_compute_subnetwork.public_subnet.id

    access_config {
      nat_ip = google_compute_address.static-ip-address.address
    }
  }

  metadata_startup_script = templatefile("./startup.tftpl",
    {
      project      = var.project
    })
} 