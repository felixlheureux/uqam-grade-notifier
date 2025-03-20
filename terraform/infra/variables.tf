# GCP variables
variable "gcp_project_id" {
  description = "Google Cloud Platform (GCP) project ID"
  type        = string
}

variable "region" {
  description = "Geographical zone for the GCP VM instance"
  type        = string
}

variable "machine_type" {
  description = "Machine type for the GCP VM instance"
  type        = string
  default     = "e2-micro"
}

# Global variables
variable "project" {
  description = "Project name"
  type        = string
}

variable "environment" {
  description = "Environment"
  type        = string
} 