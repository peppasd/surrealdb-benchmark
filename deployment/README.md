# Deployment

This directory contains the deployment intructions for the benchmark

## Prerequisites

- [gcloud CLI](https://cloud.google.com/sdk/gcloud) (462.0.1 at the time of writing)
- [terraform](https://developer.hashicorp.com/terraform/install) (v1.7.2 at the time of writing)

## Configure Terraform / GCP

- Authenticate with gcloud using `gcloud init` and choose the target project for the deployment. For the region and zone, I chose `us-central1` and `us-central1-a` respectively. This has to also be reflected in the `main.tf` file.
- Activate Application Default Credentials (ADC) using `gcloud auth application-default login`
- In the `terraform` directory run `terraform init` to initialize the terraform project
