# reloads nginx inside container without stopping by executing "nginx -s reload" inside
# docker exec web nginx -s reload

# reload container
./stop_nginx.sh
./start_nginx.sh
