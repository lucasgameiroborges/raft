VERSION=0

build:
	go build -race . 

build-release:
	go build .

test-all:
	go test -v -race ./...

test:
	go test -timeout=3s $(PACKAGE) 

gen-pb: pb/*.proto
	protoc -I pb/ pb/raft.proto --go_out=plugins=grpc:pb

setup-local-cluster:
	./scripts/setup_cluster_dir.rb

run-local-cluster: build
	./scripts/setup_cluster_dir.rb --run-cluster

build-docker-container: gen-pb
	docker build -t raft:local .

push-to-registry: build-docker-container
	docker login
	docker tag raft:local su225/raft:$(VERSION)
	docker push su225/raft:$(VERSION)

clean:
	rm -rf raft
	rm -rf local-cluster