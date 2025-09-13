cd ./front/
make dependencies
make build
cd ../
make dependencies
make target=judge build
make target=manager build