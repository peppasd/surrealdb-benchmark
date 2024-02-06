minutes=$1
threads=$2

echo "Starting benchmark with $minutes minutes per phase and $threads threads"

echo "Getting SUT internal IP..."
sut_ip="$(gcloud compute instances describe surrealdb --zone='us-central1-c' --format='get(networkInterfaces[0].networkIP)')"
echo "SUT internal IP is" $sut_ip

cmd="lg -minutes $minutes -threads $threads -url $sut_ip:8000"

echo "Running benchmark with command: $cmd"
gcloud compute ssh load-generator --ssh-flag="-ServerAliveInterval=300" --zone us-central1-c -- $cmd
echo "Benchmark finished"
