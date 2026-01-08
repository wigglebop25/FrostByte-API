# FrostByte API

**FrostByte API** is a high-performance, minimalist backend API built in **Go**, designed for blistering speed and reliability in cold-weather data processing scenarios. Built with the Go philosophy of simplicity and efficiency, it prioritizes a security-first architecture.

## ❄️ Key Features

*   **Speed & Efficiency:** Built with **Go** (Golang) and the lightweight **Chi** router for minimal overhead.
*   **Security First:** Robust **JWT** authentication with Argon2 password hashing and Role-Based Access Control (RBAC).
*   **Real-Time Updates:** Integrated **WebSockets** for live transaction and order monitoring.
*   **Architecture:** Clean **Handler-Service-Repository** pattern ensures maintainability and scalability.
*   **Data Integrity:** **MySQL** database with **GORM** for reliable data management and transactional safety.
*   **Container Ready:** Dockerized with a multi-stage build for small, secure production images.

## 🛠 Tech Stack

*   **Language:** Go 1.23+
*   **Web Framework:** Chi
*   **Database:** MySQL 8.0
*   **ORM:** GORM
*   **Real-time:** Gorilla WebSocket
*   **Authentication:** JWT (JSON Web Tokens) + Argon2

## 🚀 Getting Started

### Prerequisites
*   Docker & Docker Compose

### Prerequisites & Setup

Since Go and Docker may not be installed in your environment, follow these steps:

1.  **Install Go 1.23+**
    *   Once installed, run the following command in the project root to generate the `go.sum` file:
        ```bash
        go mod tidy
        ```

2.  **Install Docker**
    *   Once installed, run the following command to start the API and MySQL database:
        ```bash
        docker-compose up --build
        ```

3.  **Access the API**
    *   The API will be available at `http://localhost:3000` (mapped from container port 8080).

## 🧪 Testing

We use **Bruno** for API testing. A complete collection is included in this guidance package under the `bruno/` folder.

1.  Install [Bruno](https://www.usebruno.com/).
2.  Open the `guidance/bruno` collection folder in the Bruno app.
3.  The `go-api` folder contains the specific tests for this new implementation.

## 📂 Project Structure

```text
frostbyte-api/
├── cmd/
│   └── api/          # Application entry point
├── internal/
│   ├── config/       # Configuration loading
│   ├── database/     # DB connection & migrations
│   ├── domain/       # Data models (structs)
│   ├── handlers/     # HTTP Controllers
│   ├── repository/   # Data access layer
│   ├── service/      # Business logic
│   └── websocket/    # Real-time hub
├── rust_legacy/      # Archived original Rust code
├── Dockerfile
└── compose.yml
```