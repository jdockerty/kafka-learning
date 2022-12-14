.PHONY: init init-kafka init-k3d get-kafka-info
SHELL := /bin/bash

PORT := $(shell kubectl get svc -n kafka my-cluster-kafka-external-bootstrap -o jsonpath='{.spec.ports[0].nodePort}{"\n"}')
ADDRS := $(shell kubectl get nodes --output=jsonpath='{range .items[*]}{.status.addresses[?(@.type=="InternalIP")].address}{"\n"}{end}' | awk '{print $$1}')
SINGLE_ADDR :=  $(shell echo "$(ADDRS)" | awk '{print $$1}')

init-kafka:
	kubectl create -f manifests/namespace.yml
	kubectl create -f manifests/strimzi-operator.yml -n kafka
	@echo "Waiting for strimzi operator to be ready"
	sleep 15
	kubectl get pods --no-headers -n kafka | awk '{print $$1}' | head -n 1 | xargs -I % /bin/sh -c 'kubectl -n kafka wait pod % --for=condition=Ready --timeout=300s'
	kubectl create -f manifests/kafka-cluster.yml -n kafka
	kubectl create -f manifests/kafka-topic.yml -n kafka
	@echo "Kafka cluster created, please wait a couple of minutes for all resources to be fully ready."

init-k3d:
	k3d cluster create --agents 3 kafka-learning
	@echo "Verify nodes..."
	kubectl get nodes

get-kafka-info:
	@echo "NOTE: This command should only be used when all resources are available."
	@echo Port: $(PORT)
	@echo IPs: $(ADDRS)
	@echo "Use any of the IPs and port and pass it as a broker to the consumer/producer."

run-producer-local:
	@echo "Running producer!"
	go run cmd/producer/producer.go --topic test-topic --brokers "$(SINGLE_ADDR):$(PORT)"

run-consumer-local:
	@echo "Running consumer! Reading from latest offset..."
	go run cmd/consumer/consumer.go --topic test-topic --brokers "$(SINGLE_ADDR):$(PORT)"

init: init-k3d init-kafka

