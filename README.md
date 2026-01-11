# ❄️ FrostByte API

**FrostByte API** is a production-ready, high-performance backend built in **Go**. Designed for security and scalability, it features advanced Role-Based Access Control (RBAC), real-time WebSocket updates, and a hardened security architecture suitable for enterprise deployment.

## 🚀 Key Features

*   **Security First:** HTTPS/TLS, JWT Auth (Argon2), RBAC, Rate Limiting, and Secured Docker Containers.
*   **Real-Time:** Integrated WebSockets for live updates.
*   **Data Integrity:** MySQL 8.0 with GORM (Transactional).
*   **DevOps Ready:** GitHub Actions CI/CD pipeline included.

## 🛠 Tech Stack

*   **Language:** Go 1.23+
*   **Web Framework:** Chi v5
*   **Database:** MySQL 8.0
*   **Real-time:** Gorilla WebSocket
*   **Deployment:** Docker, Docker Compose, GitHub Actions

## 🏁 Getting Started

### 1. Prerequisites
*   Docker & Docker Compose
*   SSL Certificates (`server.crt`, `server.key`) in the root.
*   `.env` file (copy from `.env.example`).

### 2. Quick Start
```bash
# Start the stack (HTTPS enabled)
docker-compose up -d --build
```
Access at: **`https://127.0.0.1:3000`**

### 3. Reset Database
```bash
docker-compose down -v
docker-compose up -d --build
```

## 🧪 API Testing (Bruno)

Import the `bruno/` folder into [Bruno](https://www.usebruno.com/).
*   **00_Auth**: Login/Register (First user registered becomes **Admin**).
*   **01_Admin_Workflow**: Full system control.
*   **02_Cashier_Workflow**: Order management.
*   **03_Customer_Workflow**: Shopping & History.

## 🚢 Deployment

This project includes a **GitHub Actions** workflow (`.github/workflows/deploy.yml`) that automatically deploys to a VPS.

**Requirements:**
1.  **Docker Hub Secrets:** `DOCKER_USERNAME`, `DOCKER_PASSWORD`.
2.  **Server Secrets:** `HOST` (IP), `USERNAME`, `SSH_KEY`.
3.  **Server Setup:** Install Docker and create `.env` file manually on the server.