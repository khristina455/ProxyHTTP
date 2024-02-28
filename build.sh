cd ./scripts
chmod +x ./gen_ca.sh
./gen_ca.sh
cd ../..
docker compose build
docker compose up -d
