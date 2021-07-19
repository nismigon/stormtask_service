cd ..
docker build -t stormtask_lint_image -f Dockerfile.server.lint .
docker run --rm stormtask_lint_image 