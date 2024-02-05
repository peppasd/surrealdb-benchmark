terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "5.15.0"
    }
  }
}

provider "google" {
  project = "csb-ws2324"
  region  = "us-central1"
  zone    = "us-central1-c"
}

### NETWORK
resource "google_compute_network" "vpc_network" {
  name                    = "benchmark-network"
  auto_create_subnetworks = true
}

### FIREWALL
resource "google_compute_firewall" "all" {
  name = "allow-all"
  allow {
    protocol = "tcp"
    ports    = ["0-65535"]
  }
  network       = google_compute_network.vpc_network.id
  source_ranges = ["0.0.0.0/0"]
}

### SUT INSTANCE
resource "google_compute_instance" "sut" {
  name         = "surrealdb"
  machine_type = "n2d-standard-2"

  boot_disk {
    initialize_params {
      size  = 20
      type  = "pd-ssd"
      image = "ubuntu-2204-jammy-v20240126"
    }
  }

  metadata_startup_script = file("startup_sut.sh")

  scratch_disk {
    interface = "NVME"
  }

  network_interface {
    network = google_compute_network.vpc_network.id
    access_config {
    }
  }
}

### LOAD GENERATOR INSTANCE
resource "google_compute_instance" "lg" {
  name         = "load-generator"
  machine_type = "n2d-standard-2"

  boot_disk {
    initialize_params {
      size  = 20
      type  = "pd-ssd"
      image = "ubuntu-2204-jammy-v20240126"
    }
  }

  metadata_startup_script = file("startup_lg.sh")

  scratch_disk {
    interface = "NVME"
  }

  network_interface {
    network = google_compute_network.vpc_network.id
    access_config {
    }
  }
}
