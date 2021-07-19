cd server
docker build -t stormtask_image_test -f Dockerfile.test .
docker-compose -f docker/compose_test.yaml -p stormtask_compose_test up -d