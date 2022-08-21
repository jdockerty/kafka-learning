# Kafka Learning

The `manifests` folder can be used to deploy the [`strimzi`](https://strimzi.io/) Kafka operator into [`k3d`](https://k3d.io/), by running `make init`.


Once the operator is deployed, you can produce some fake data into the pre-created `test-topic` by running

    # Run these in separate terminals as they are blocking calls
    make run-producer-local
    make run-consumer-local
