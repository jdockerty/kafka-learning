# Kafka Learning

The `manifests` folder can be used to deploy the [`strimzi`](https://strimzi.io/) Kafka operator into [`k3d`](https://k3d.io/), by running `make init`.


Once the operator is deployed, you can produce some fake data into the pre-created `test-topic` by running

    # Run these in separate terminals as they are blocking calls
    make run-producer-local
    make run-consumer-local

This assumes the use of Strimzi, if using Confluent Cloud it requires running the `go` binary directly, as an API and secret key must be passed for authentication

    go run cmd/producer/producer.go --topic test-topic --brokers "<bootstrap_server>" --api-key "<api_key>" --secret-key "<secret_key>"
    go run cmd/consumer/consumer.go --topic test-topic --brokers "<bootstrap_server>" --api_key "<api_key>" --secret-key "<secret_key>"
