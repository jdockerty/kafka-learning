.PHONY: init init-kafka init-k3d

init-kafka:
	kubectl create -f manifests/namespace.yml
	kubectl create -f manifests/strimzi-quickstart.yml -n kafka
	echo "Waiting for strimzi operator to be ready"
	kubectl get pods --no-headers -n kafka | awk '{print $1}' | head -n 1 | xargs -I % /bin/sh -c 'kubectl wait --for condition=Ready --timeout=300s pod/%'
	kubectl create -f manifests/strimzi-kafka-single.yml -n kafka

init-k3d:
	k3d cluster create --agents 3 kafka-learning
	kubectl get nodes # Verify

init: init-k3d init-kafka

