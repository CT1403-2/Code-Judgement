cd ./front/
make dependencies
make build
cd ../
make dependencies
make target=judge build
make target=manager build

docker build -f ./manifests/runner.Dockerfile -t runner:v0.0.7 .
