VERSION=v1.0.2
DOCKER_NW=dockerraft
CLUSTER_SZ=3

RAFT_HEADLESS_SVC=raft-headless-svc.yaml
K8S_DEPLOY_DESCRIPTOR=raft-k8s-deploy.yaml
PERSISTENT_VOLUMES=volumes.yaml

build:
	go build -race . 

build-release:
	go build .

test-all:
	go test -v -race ./...

test:
	go test -timeout=3s $(PACKAGE) 

setup-local-cluster:
	./scripts/cluster_manager.rb --generate --dry-run

run-local-cluster: build
	./scripts/cluster_manager.rb --generate --launch \
		--cluster-size=$(CLUSTER_SZ)

run-docker-cluster:
	./scripts/cluster_manager.rb --generate --launch \
		--cluster-size=$(CLUSTER_SZ) \
		--docker-mode --docker-network-name=$(DOCKER_NW)

destroy-docker-cluster:
	./scripts/cluster_manager.rb --docker-mode --docker-destroy \
		--cluster-size=$(CLUSTER_SZ)

build-docker-container:
	docker build -t raft:local .

push-to-registry: build-docker-container
	docker login
	docker tag raft:local lucasgameiroborges/raft:$(VERSION)
	docker push lucasgameiroborges/raft:$(VERSION)

k8s-deploy:
	microk8s kubectl create -f $(PERSISTENT_VOLUMES)
	microk8s kubectl create -f $(RAFT_HEADLESS_SVC)
	microk8s kubectl create -f $(K8S_DEPLOY_DESCRIPTOR)

k8s-check:
	microk8s kubectl get pvc -n raft-k8s
	microk8s kubectl get pv -n raft-k8s
	microk8s kubectl get pods -n raft-k8s
	microk8s kubectl get services -n raft-k8s
	microk8s kubectl get statefulset -n raft-k8s

k8s-undeploy:
	microk8s kubectl delete statefulset --all -n raft-k8s
	microk8s kubectl delete pod --all -n raft-k8s
	microk8s kubectl delete pvc --all -n raft-k8s
	microk8s kubectl delete pv --all -n raft-k8s
	microk8s kubectl delete service --all -n raft-k8s

k8s-describe:
	microk8s kubectl describe pod raft000-0 -n raft-k8s


mini-deploy:
	minikube kubectl -- create -f $(PERSISTENT_VOLUMES)
	minikube kubectl -- create -f $(RAFT_HEADLESS_SVC)
	minikube kubectl -- create -f $(K8S_DEPLOY_DESCRIPTOR)

mini-check:
	minikube kubectl -- get pvc -n raft-k8s
	minikube kubectl -- get pv -n raft-k8s
	minikube kubectl -- get pods -n raft-k8s
	minikube kubectl -- get services -n raft-k8s
	minikube kubectl -- get statefulset -n raft-k8s

mini-undeploy:
	minikube kubectl -- delete statefulset --all -n raft-k8s
	minikube kubectl -- delete pod --all -n raft-k8s
	minikube kubectl -- delete pvc --all -n raft-k8s
	minikube kubectl -- delete pv --all -n raft-k8s
	minikube kubectl -- delete service --all -n raft-k8s

mini-describe:
	minikube kubectl -- describe pod raft000-1 -n raft-k8s

clean:
	rm -rf raft
	rm -rf local-cluster