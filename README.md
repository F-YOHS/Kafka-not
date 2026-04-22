# Kafka-not

Event-driven microservices demo built with Go and Apache Kafka.

A user registration event triggers a notification pipeline — email, SMS, and analytics services react to events independently via Kafka topics.

## Architecture

```
User Service (producer)
    │
    ▼
Kafka (4 topics)
    ├── user-events        → Analytics Service
    ├── payment-events     → Analytics Service
    ├── notification-events → Email Service
    │                       → SMS Service
    └── analytics-events   (produced by Analytics)
```

## Services

| Service           | Port | Role                                      |
|-------------------|------|-------------------------------------------|
| user-service      | 8080 | Accepts HTTP requests, produces events    |
| email-service     | —    | Consumes notification-events, sends email |
| sms-service       | —    | Consumes notification-events, sends SMS   |
| analytics-service | —    | Aggregates user + payment events          |
| Kafdrop (UI)      | 9000 | Kafka cluster monitoring                  |

## Tech Stack

- **Go** — all microservices
- **Apache Kafka** — message broker (Confluent 7.3.0)
- **Zookeeper** — Kafka coordination
- **Docker / Docker Compose** — containerization
- **Kafdrop** — Kafka web UI

## Getting Started

**Prerequisites:** Docker, Docker Compose, Make

```bash
# Clone the repo
git clone https://github.com/F-YOHS/Kafka-not.git
cd Kafka-not

# Copy environment config
cp .env.example .env

# Build and start all services
make build
make up
```

Kafka UI available at: http://localhost:9000

## Usage

Send a test registration event:

```bash
make request
```

This POSTs a sample user (email + phone) to the user-service, which publishes to Kafka. Email, SMS, and analytics services consume the event automatically.

## Commands

| Command        | Description                        |
|----------------|------------------------------------|
| `make build`   | Build all Docker images            |
| `make up`      | Start all containers               |
| `make down`    | Stop and remove containers         |
| `make logs`    | Tail logs from all services        |
| `make request` | Send a test registration request   |

## Topics

| Topic               | Partitions | Producers         | Consumers                  |
|---------------------|------------|-------------------|----------------------------|
| user-events         | 3          | user-service      | analytics-service          |
| payment-events      | 3          | —                 | analytics-service          |
| notification-events | 3          | user-service      | email-service, sms-service |
| analytics-events    | 3          | analytics-service | —                          |
