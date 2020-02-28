MODULE=kannan.ieee.org
REGISTRY=kvaradha
VERSIONS=0.0.1
GOARGS=-race
GOLINUX_ENV=CGO_ENABLED=0 GOOS=linux

CONTAINERS=hello-world-server hello-world-client
LINUX_OBJECTS=server/server.Linux client/client.Linux
NATIVE_OBJECTS=server/server client/client
SERVER_SOURCES=server/server.go proto/helloWorld/helloWorld.pb.go
CLIENT_SOURCES=client/client.go proto/helloWorld/helloWorld.pb.go

all: $(CONTAINERS) $(NATIVE_OBJECTS)

deploy: $(CONTAINERS) server/k8s-d-server.yml client/k8s-d-client.yml
	docker push kvaradha/hello-world-server:$(VERSIONS)
	docker push kvaradha/hello-world-server:latest
	docker push kvaradha/hello-world-client:$(VERSIONS)
	docker push kvaradha/hello-world-client:latest
	-kubectl create -f server/k8s-d-server.yml
	-kubectl create -f client/k8s-d-client.yml

hello-world-server: server/Dockerfile server/server.Linux
	docker build -t $(REGISTRY)/$@:$(VERSIONS) server
	docker tag $(REGISTRY)/$@:$(VERSIONS) $(REGISTRY)/$@:latest

hello-world-client: client/Dockerfile client/client.Linux
	docker build -t $(REGISTRY)/$@:$(VERSIONS) client
	docker tag $(REGISTRY)/$@:$(VERSIONS) $(REGISTRY)/$@:latest

server/server.Linux: $(SERVER_SOURCES)
	$(GOLINUX_ENV) go build -a -installsuffix cgo -o $@ ./server

server/server: $(SERVER_SOURCES)
	go build $(GOARGS) -o $@ ./server

client/client.Linux: $(CLIENT_SOURCES)
	$(GOLINUX_ENV) go build -a -installsuffix cgo -o $@ ./client

client/client: $(CLIENT_SOURCES)
	go build $(GOARGS) -o $@ ./client

proto/helloWorld/helloWorld.pb.go: proto/helloWorld/helloWorld.proto
	protoc -I ./proto/helloWorld --go_out=plugins=grpc:./proto/helloWorld helloWorld.proto

modules: 
	@rm -f go.mod go.sum
	go mod init $(MODULE)

clean:
	@rm -f $(LINUX_OBJECTS) $(NATIVE_OBJECTS)

spotless: clean
	@rm -f proto/helloWorld/helloWorld.pb.go
	@rm -f go.mod go.sum

