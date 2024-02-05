cmd="surreal start --allow-net file://surreal.db";

export DEBIAN_FRONTEND=noninteractive
sudo apt update
sudo apt upgrade -y
curl -sSf https://install.surrealdb.com | sh
eval "${cmd}" &>/dev/null & disown;

