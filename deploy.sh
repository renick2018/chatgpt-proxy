cd chatgpt-proxy
git pull
cd ../chatgpt-api-server
git checkout api-server-3.5
git pull
cd ../
cp chatgpt-proxy/docker-compose.yml .
docker compose -f docker-compose.yml up --build -d
docker rmi `docker images -aq`
