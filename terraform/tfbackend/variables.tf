# GCP variables
variable "gcp_project_id" {
  description = "Google Cloud Platform (GCP) project ID"
  type        = string
}

variable "region" {
  description = "Geographical zone for the GCP VM instance"
  type        = string
} 