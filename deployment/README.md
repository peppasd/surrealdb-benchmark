# Deployment

This directory contains the deployment intructions for the benchmark

## Prerequisites

- [gcloud CLI](https://cloud.google.com/sdk/gcloud) (462.0.1 at the time of writing)
- [terraform](https://developer.hashicorp.com/terraform/install) (v1.7.2 at the time of writing)

## Configure Terraform / GCP

- Authenticate with gcloud using `gcloud init` and choose the target project for the deployment. For the region and zone, I chose `us-central1` and `us-central1-c` respectively. In case you change that, it has to also be reflected in the `main.tf` file and all the scripts.
- Activate Application Default Credentials (ADC) using `gcloud auth application-default login`
- In the `terraform` directory run `terraform init` to initialize the terraform project
- In the `terraform` directory run `terraform apply` to deploy the infrastructure

## Run the benchmark

There is nothing else to configure to run the benchmark. The startup script will take care of setting up the environment and loading the necessary data. However this process can take a while after the deployment is finished. To check if both instances are ready for the benchmark, you can run the following command:

```bash
bash check_readiness.sh
```

If both instances return `Ready!`, you can run the benchmark using the following command:

```bash
bash run_benchmark.sh <minutes_per_phase> <number_of_threads>
```

Where `<minutes_per_phase>` is the number of minutes you want to run the benchmark for each phase and and `<number_of_threads>` is the number of threads you want to use for the benchmark. There are 3 phases in total: `REST`, `Websocket` and `SDK`. Example:

```bash
bash run_benchmark.sh 20 3
```

This will run the benchmark for 20 minutes for each phase using 3 threads.

The script will run the benchmark in the background. If you want to see the logs, you can run the following command:

```bash
bash check_logs.sh
```

After the benchmark is finished, you can download the results using the following command:

```bash
bash download_results.sh
```

To remove the old logs, you can run the following command:

```bash
bash delete_logs.sh
```

## Rerun the benchmark

You can rerun the benchmark as many consecutive times as you want, without any extra configuration, because the load generator
will leave the data in the database at the same state as before. In the improbable case that the load generator fails or if you stop the banchmark manually, you have to redeploy the infrastructure and run the benchmark again in order to ensure that the database is in the correct state.

IMPORTANT: Any rerun will overwrite the previous results, so make sure to download the existing results before running the benchmark again.

## Destroy the infrastructure

After you are done with the benchmark, you can destroy the infrastructure using the following command in the `terraform` directory:

```bash
terraform apply -destroy
```
