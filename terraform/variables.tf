variable "project_id" {
  description = "ID du projet GCP"
  type        = string
}

variable "region" {
  description = "Région GCP"
  type        = string
  default     = "us-central1"
}

variable "zone" {
  description = "Zone GCP"
  type        = string
  default     = "us-central1-a"
}

variable "ssh_user" {
  description = "Nom d'utilisateur SSH"
  type        = string
}

variable "ssh_pub_key_path" {
  description = "Chemin vers la clé publique SSH"
  type        = string
} 