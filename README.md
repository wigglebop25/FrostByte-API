# FrostByte Full-Stack System

FrostByte is a production-ready, high-performance restaurant management system featuring a Go (Chi) Backend, a SvelteKit (TypeScript) Web Frontend, and a secure Caddy Reverse Proxy. Designed for speed and security, it includes real-time WebSocket updates, advanced RBAC, and automated Azure deployment.

## Live Demo
* **Web Application**: [https://frostbyte-api.southeastasia.cloudapp.azure.com](https://frostbyte-api.southeastasia.cloudapp.azure.com)
* **API Version**: v1 (Base URL: /api/v1)
* **Status**: Deployed & Production-Ready

---

## System Architecture

### 1. Web Frontend (/frontend)
* **Framework**: SvelteKit (Svelte 5) with TypeScript.
* **Styling**: Tailwind CSS with a modern Glassmorphism design.
* **Features**: Responsive Dashboard, Real-time Order Tracking, Admin Management, and Dark Mode support.

### 2. Mobile Client (Android)
* **Repository**: [wigglebop25/RestaurantClient](https://github.com/wigglebop25/RestaurantClient)
* **Framework**: Native Android (Kotlin) with Jetpack Compose.
* **Integration**: Leverages the Go API for high-performance data fetching and WebSocket integration for real-time order status notifications.

### 3. Go Backend (/cmd/api)
* **Engine**: Go 1.23+ using the Chi router.
* **Database**: MySQL 8.0 with GORM (Transactional integrity).
* **Security**: JWT Authentication (Argon2), Role-Based Access Control (Admin/Cashier/Customer), and Rate Limiting.
* **Real-time**: High-concurrency WebSocket Hub for live order notifications.

### 4. Infrastructure & DevOps
* **Reverse Proxy**: Caddy handles automatic HTTPS (Let's Encrypt), load balancing, and static asset serving.
* **Containerization**: Fully Dockerized (Non-root, multi-stage builds).
* **Cloud**: Hosted on Azure VM (Southeast Asia).
* **CI/CD**: Automated GitHub Actions pipeline for building and deploying.

---

## Tech Stack

| Layer | Technologies |
| :--- | :--- |
| **Frontend** | SvelteKit, TypeScript, Tailwind CSS, Lucide, ApexCharts |
| **Mobile** | Android SDK, Kotlin, Jetpack Compose, Retrofit |
| **Backend** | Go, Chi, GORM, Gorilla WebSocket |
| **Database** | MySQL 8.0 |
| **Proxy/SSL** | Caddy v2 (Auto-HTTPS) |
| **DevOps** | Docker, GitHub Actions, Azure VM |

---

## Key Documentation

### Client Applications
* [**Mobile Client Repository**](https://github.com/wigglebop25/RestaurantClient) - Native Android implementation.
* [**Frontend Setup**](frontend/README.md) - How to run the web app locally.
* [Web Integration Guide](frontend/WEB_INTEGRATION_GUIDE.md) - Svelte & Styling details.
* [API Communication Guide](API_COMMUNICATION_GUIDE.md) - How clients talk to the API.

### Backend & API
* [**API Reference**](API_REFERENCE_COMPLETE.md) - Complete list of endpoints and payloads.
* [WebSocket Guide](WEBSOCKET_GUIDE.md) - Real-time event documentation.
* [Android Backend Updates](ANDROID_BACKEND_UPDATES.md) - Latest API changes for mobile.

### Deployment & DevOps
* [Azure Deployment](docs/AZURE_DEPLOY.md) - Instructions for the cloud infrastructure.
* [Master Plan](docs/MASTER_PLAN.md) - Original project architecture and goals.

---
Copyright 2026 FrostByte System | [LICENSE](LICENSE)
