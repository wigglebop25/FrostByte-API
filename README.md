# ❄️ FrostByte API

**FrostByte API** is a production-ready, high-performance backend built in **Go**. Designed for security and scalability, it features advanced Role-Based Access Control (RBAC), real-time WebSocket updates, and a hardened security architecture suitable for enterprise deployment.

## 🚀 Key Features

*   **Security First:**
    *   **HTTPS/TLS Enforcement:** Strict Transport Security (HSTS) and secure headers.
    *   **Authentication:** JWT (JSON Web Tokens) with Argon2 hashing.
    *   **RBAC:** Granular control (Admin, Cashier, Customer workflows).
    *   **Rate Limiting:** Built-in protection against abuse (Token Bucket).
    *   **Container Security:** Non-root Docker images with no baked-in secrets.
*   **Real-Time Updates:** Integrated **WebSockets** for live order status monitoring.
*   **Architecture:** Clean **Handler-Service-Repository** pattern.
*   **Data Integrity:** **MySQL** 8.0 with GORM, fully transactional.
*   **DevOps Ready:** Includes **CI/CD** workflows (GitHub Actions) and Docker Compose.

## 🛠 Tech Stack

*   **Language:** Go 1.23+
*   **Web Framework:** Chi v5
*   **Database:** MySQL 8.0
*   **Real-time:** Gorilla WebSocket
*   **Deployment:** Docker, Docker Compose, GitHub Actions

## 🏁 Getting Started

### 1. Prerequisites
*   Docker & Docker Compose
*   A pair of SSL certificates (`server.crt` and `server.key`) in the root directory.

### 2. Quick Start (Docker)
The easiest way to run the API is using the included production-ready Docker setup.

```bash
# 1. Clone the repository
git clone https://github.com/your-repo/frostbyte-api.git
cd frostbyte-api

# 2. Configure Environment
# Copy the example and edit strictly (Do NOT commit .env)
cp .env.example .env

# 3. Start the Stack (HTTPS enabled)
docker-compose up -d --build
```

The API will be available at **`https://127.0.0.1:3000`**.

> **Note:** Since we use self-signed certificates for development, you may need to accept the security warning in your browser/client.

### 3. Resetting the Database
To clear all data and start fresh (useful for testing "First User = Admin" logic):
```bash
docker-compose down -v
docker-compose up -d --build
```

## 🧪 API Testing (Bruno)

We use **Bruno** for API testing. The `bruno/` folder is organized by user roles to simplify workflow testing.

1.  **Install [Bruno](https://www.usebruno.com/)**.
2.  **Open Collection**: Point Bruno to the `bruno/` folder in this repo.
3.  **Workflows**:
    *   `00_Auth`: Register & Login (Start here).
    *   `01_Admin_Workflow`: Full system control (Users, Roles, Products, Analytics).
    *   `02_Cashier_Workflow`: Manage Orders.
    *   `03_Customer_Workflow`: Shop & View Order History.

## 🚢 Deployment

For detailed production deployment instructions, including setting up a VPS and GitHub Secrets, see **[DEPLOYMENT.md](DEPLOYMENT.md)**.

## 📂 Project Structure

```text
frostbyte-api/
├── .github/workflows # CI/CD Pipelines
├── bruno/            # API Request Collections
├── cmd/api/          # Application entry point
├── internal/
│   ├── config/       # Configuration
│   ├── database/     # DB connection & migrations
│   ├── domain/       # Data models (structs)
│   ├── handlers/     # HTTP Controllers & Middleware
│   ├── repository/   # Data access layer
│   ├── service/      # Business logic
│   └── websocket/    # Real-time hub
├── Dockerfile        # Secured multi-stage build
└── compose.yml       # Production-ready compose file
```
