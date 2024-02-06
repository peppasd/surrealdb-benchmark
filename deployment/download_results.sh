echo "Downloading results..."
timestamp=$(date +%s)
gcloud compute scp load-generator:results.sqlite ./results-$timestamp.sqlite --zone='us-central1-c'
echo "Results downloaded to ./results-$timestamp.sqlite"