cmd="sudo surreal start file://surreal.db";

export DEBIAN_FRONTEND=noninteractive
sudo apt update
sudo apt upgrade -y
curl -sSf https://install.surrealdb.com | sh
eval "${cmd}" &>/dev/null & disown;

curl -L https://github.com/peppasd/surrealdb-benchmark/raw/main/prepare_db/surrealdata.tar.gz | tar -xz
sudo surreal import --conn http://localhost:8000 --ns benchmark --db benchmark db.surql
sudo rm db.surql

touch /done
