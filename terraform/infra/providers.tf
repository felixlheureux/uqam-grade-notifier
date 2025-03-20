terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = ">= 5.33.0"
    }
    random = {
      source = "hashicorp/random"
      version = ">= 3.6.2"
    }
  }
  required_version = ">= 1.8.5"

  backend "gcs" {
    bucket  = "a7d1444320ed4f61-bucket-tfstate"
    prefix  = "terraform/state"
  }
}

# Providers
provider "google" {
  project               = var.gcp_project_id
  region                = var.region
  zone                  = "${var.region}-b"
  user_project_override = true
  billing_project       = var.gcp_project_id
}
provider "random" {
} 