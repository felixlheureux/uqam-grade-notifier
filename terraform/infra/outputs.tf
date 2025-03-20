output "instance_static_ip_address" {
  value = google_compute_address.static-ip-address.address
} 