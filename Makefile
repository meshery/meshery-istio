protoc-setup:
	cd meshes
	wget https://raw.githubusercontent.com/layer5io/meshery/master/meshes/meshops.proto

proto:	
	protoc -I meshes/ meshes/meshops.proto --go_out=plugins=grpc:./meshes/

docker:
	docker build -t layer5/meshery-istio .

docker-run:
	(docker rm -f meshery-istio) || true
	docker run --name meshery-istio -d \
	-p 10000:10000 \
	-e DEBUG=true \
	layer5/meshery-istio

run:
	DEBUG=true GOPROXY=direct GOSUMDB=off go run main.go