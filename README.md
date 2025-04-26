Real-Time Messaging Application
An event-driven, real-time messaging application based on microservices architecture, built with Go, PostgreSQL, and Apache Kafka.
Architecture Overview
This application follows a microservices architecture with the following components:

WebSocket/SSE Service: Handles real-time communication with clients, consuming from and publishing to Kafka
Persistence Service: Consumes messages from Kafka and stores them in PostgreSQL
Authentication Service: Provides secure, stateless authentication for all services
Kafka: Message broker for asynchronous communication between services
PostgreSQL: Persistent storage for messages and user data

Prerequisites

Go 1.20+
PostgreSQL 14+
Apache Kafka 3.0+
Docker and Docker Compose
Kubernetes cluster (for production deployment)
Helm (optional, for Kubernetes deployment)

Getting Started
Local Development Setup

Clone the repository:
bashgit clone https://github.com/yourusername/messaging-app.git
cd messaging-app

Set up environment variables:
bashcp .env.example .env
Edit the .env file with your configuration settings.
Start the dependencies (PostgreSQL and Kafka) using Docker Compose:
bashdocker-compose up -d postgres kafka

Install Go dependencies:
bashgo mod download

Run each service individually for development:
bash# Terminal 1: WebSocket/SSE Service
cd services/websocket
go run main.go

# Terminal 2: Persistence Service
cd services/persistence
go run main.go

# Terminal 3: Authentication Service
cd services/auth
go run main.go


Running with Docker Compose
To run the entire application stack using Docker Compose:
bashdocker-compose up -d
This will start all services, PostgreSQL, and Kafka in separate containers.