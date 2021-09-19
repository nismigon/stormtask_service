docker build -t stormtask_dev_image .
docker-compose -f docker/compose_dev.yaml -p stormtask_compose_dev up -d