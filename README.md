# FrostByte API

FrostByte API is a production-ready, high-performance backend built in Go. Designed for security and scalability, it features advanced Role-Based Access Control (RBAC), real-time WebSocket updates, and a hardened security architecture suitable for enterprise deployment.

## Key Features

* **Security First**: HTTPS/TLS, JWT Auth (Argon2), RBAC, Rate Limiting, and Secured Docker Containers.
* **Real-Time**: Integrated WebSockets for live updates.
* **Data Integrity**: MySQL 8.0 with GORM (Transactional).
* **DevOps Ready**: GitHub Actions CI/CD pipeline included.
* **Mobile-Ready**: 
 * **Server-Side Analytics**: Pre-calculated dashboard stats for performance.
 * **Dual-Token Auth**: Secure Access (15m) & Refresh (7d) token flow.
 * **Optimized Data**: Pagination, Fuzzy Search, and Token Bucket Rate Limiting.

## Tech Stack

* **Language**: Go 1.23+
* **Web Framework**: Chi v5
* **Database**: MySQL 8.0
* **Real-time**: Gorilla WebSocket
* **Deployment**: Docker, Docker Compose, GitHub Actions

## Documentation

* [**Latest Updates (Android Integration)**](ANDROID_BACKEND_UPDATES.md) - **Read this for new API changes.**
* [Implementation Guide](docs/IMPLEMENTATION_GUIDE.md) - Step-by-step guide on how the API was built.
* [Azure Deployment Guide](docs/AZURE_DEPLOY.md) - Instructions for deploying to Azure.
* [Master Plan](docs/MASTER_PLAN.md) - Original project plan and architecture.
* [Deployment Details](docs/DEPLOYMENT.md) - Additional deployment information.
