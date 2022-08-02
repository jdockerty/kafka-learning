
init:
	kubectl create -f manifests/ -n kafka

init-k3d:
	k3d cluster create --agents 3 kafka-learning
	kubectl get nodes # Verify
