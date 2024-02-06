#Autostart services https://askubuntu.com/questions/1367139/apt-get-upgrade-auto-restart-services
export DEBIAN_FRONTEND=noninteractive
sudo apt update
sudo apt upgrade -y
sudo apt install -y -q git
sudo apt install -y -q golang-go

sudo git clone https://github.com/peppasd/surrealdb-benchmark.git
cd surrealdb-benchmark/load_generator
sudo go mod download
sudo go build .
sudo cp load_generator /usr/local/bin/lg
sudo rm -rf surrealdb-benchmark

touch /done
