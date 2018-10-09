docker rm -f $(docker ps -a | grep dev- | awk '{print $1}' )
docker rmi -f $(docker images | grep dev- | awk '{print $3}')
