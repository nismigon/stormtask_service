cd server
docker build -t stormtask_image_lint -f Dockerfile.lint .
docker run stormtask_image_lint