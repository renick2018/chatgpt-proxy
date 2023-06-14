cp chatgpt-proxy/docker-compose.yml .
cd chatgpt-proxy
git pull
cd ../chatgpt-api-server
git checkout api-server-3.5
git pull
cd ../
docker compose -f docker-compose.yml up --build -d
