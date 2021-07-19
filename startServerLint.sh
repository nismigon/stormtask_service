cd server
docker build -t stormtask_lint_image -f Dockerfile.lint .
docker run stormtask_lint_image