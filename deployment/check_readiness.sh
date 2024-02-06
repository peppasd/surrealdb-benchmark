cmd='test -f /done && echo "Ready!" || echo "Not ready!"'

echo "Checking readiness of the sut"
gcloud compute ssh surrealdb --zone us-central1-c -- $cmd

echo "Checking readiness of the load generator"
gcloud compute ssh load-generator --zone us-central1-c -- $cmd
