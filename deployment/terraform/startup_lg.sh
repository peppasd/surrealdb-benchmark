export DEBIAN_FRONTEND=noninteractive
sudo apt update
sudo apt upgrade -y
sudo apt install -y -q git
sudo apt install -y -q golang-go

cd ~/
git clone https://github.com/peppasd/surrealdb-benchmark.git
cd surrealdb-benchmark/load_generator
go mod download
go build .
cp load_generator ../../load_generator