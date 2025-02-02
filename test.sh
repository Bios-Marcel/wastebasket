# /bin/sh

docker build -f test_linux.Dockerfile -t temp_test . && docker run temp_test && docker container prune --force && docker image rm temp_test
