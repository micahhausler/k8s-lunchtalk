export REDIS_PASSWORD=$(head -c 32 /dev/random  | base64 )

export REDIS_HOST="192.168.99.100"
export REDIS_PASSWORD="password"
docker run -d --name redis -p 6379:6379 redis:alpine redis-server --requirepass $REDIS_PASSWORD
