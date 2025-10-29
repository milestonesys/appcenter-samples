# Kafka Topics Sample

A sample application that demonstrates producing and consuming messages with Apache Kafka in Runtime Platform. The sample consists of two Go services: a producer that continuously publishes test messages and a consumer that reads them from the same topic.

## Overview

This sample demonstrates:
- Declaring Kafka topics in `app-definition.yaml`
- Publishing messages from a Go producer
- Consuming messages
- Using platform-provided environment variables for Kafka connectivity

## What is Apache Kafka?

Apache Kafka is an open-source distributed event streaming platform used by thousands of companies for high-performance data pipelines, streaming analytics, data integration, and mission-critical applications.[Kafka](https://kafka.apache.org/) Core concepts:
- **Topic**: A named stream of records. Topics are split into partitions for scalability.
- **Partition**: is an essential components within Kafka's distributed architecture that enable Kafka to scale horizontally, allowing for efficient parallel data processing
- **Producer**: Writes records to topics.
- **Consumer**: Reads records from topics.
- **Broker**: A Kafka server that stores partitions and serves produce / fetch requests.

## Sample Structure

```
kafka-topics/
├── app-definition.yaml          # Declares Kafka producer & consumer topics and services
├── Makefile                     # Build and deployment commands (app-builder + Docker)
├── README.md                    # This documentation
└── containers/
    └── kafka/
        ├── producer/
        │   ├── Dockerfile       # Producer image build instructions
        │   └── src/
        │       ├── go.mod       # Go module and franz-go dependency
        │       └── main.go      # Produces test messages continuously
        └── consumer/
            ├── Dockerfile       # Consumer image build instructions
            └── src/
                ├── go.mod       # Go module and franz-go dependency
                └── main.go      # Consumes messages until interrupted
```

## Application Definition Highlights

Excerpt from `app-definition.yaml`:
```yaml
messaging:
  kafka:
    consumerTopics:
    - name: "samples.my-topic"
    producerTopics:
    - name: "samples.my-topic"
```
Both services reference the same topic (`samples.my-topic`). The platform provisions access and injects connection details.

Services:
- `samples-producer-service` -> container `sandbox.io/kafka/producer:1.0.0`
- `samples-consumer-service` -> container `sandbox.io/kafka/consumer:1.0.0`

## Environment Variables

The code expects `KAFKA_BOOTSTRAP_SERVER` to be set by the platform; it is automatically injected when you define Kafka topics in the `messaging.kafka` block of your `app-definition.yaml`.

Producer snippet:
```go
bootstrapServer := os.Getenv("KAFKA_BOOTSTRAP_SERVER")
cl, err := kgo.NewClient(kgo.SeedBrokers(bootstrapServer))
```

Consumer snippet (simplified):
```go
bootstrapServer := os.Getenv("KAFKA_BOOTSTRAP_SERVER")
cl, err := kgo.NewClient(
  kgo.SeedBrokers(bootstrapServer),
  kgo.ConsumeTopics("samples.my-topic"),
)
```

## How to Build and Deploy

First log into your running cluster:
```bash
cd src/kafka-topics
make login
make build   # builds both Docker images and the chart
```

### Deploy to App Center
Push the application artifacts and install:
```bash
make push              # pushes images + chart to sandbox registry
make install-from-repo # installs the app from the repository
```

### Verify Deployment
```bash
make list    # list installed applications
make events  # view application events
```

### View Logs
Use standard Kubernetes tooling or the platform dashboard to inspect logs:
```bash
# Replace <producer-pod> / <consumer-pod> with actual pod names
kubectl logs -n kafka-app <producer-pod>
kubectl logs -n kafka-app <consumer-pod>
```
You should see "producing" lines in the producer log and "consuming" lines in the consumer log.

### Uninstall / Cleanup
```bash
make uninstall
make remove      # remove release metadata
make remove-image  # optional: remove built local images (see Makefile targets)
```
Or remove images selectively:
```bash
make remove-image  # invokes docker rmi for sample images
```