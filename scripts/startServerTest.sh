cd ..
docker build -t stormtask_test_image -f Dockerfile.server.test .
docker-compose -f docker/compose_test.yaml -p stormtask_compose_test up -d