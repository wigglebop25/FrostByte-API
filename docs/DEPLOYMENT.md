# FrostByte API Deployment Guide

This guide details how to deploy the FrostByte API using Docker.

## Prerequisites
- **Docker & Docker Compose**: Ensure these are installed on your server (e.g., Ubuntu VPS, AWS EC2, DigitalOcean Droplet).
- **Git**: To clone the repository.
- **SSL Certificates**: You need `server.crt` and `server.key` for HTTPS.

## 1. Environment Configuration
**Never commit `.env` to Git.** Create it manually on the server.

1. Copy the example:
 ```bash
 cp .env.example .env
 ```
2. Edit `.env` with production values:
 ```ini
 SERVER_PORT=443
 DB_HOST=db
 DB_PORT=3306
 DB_USER=frost_admin
 DB_PASSWORD=YOUR_STRONG_PASSWORD
 DB_NAME=frost_db
 JWT_SECRET=YOUR_VERY_LONG_COMPLEX_SECRET_KEY
 ```

## 2. Docker Deployment
We use `docker-compose` to run the API and Database.

### Steps:
1. **Clone the Repo**:
 ```bash
 git clone https://github.com/your-repo/frostbyte-api.git
 cd frostbyte-api
 ```
2. **Add Certificates**:
 Place your valid `server.crt` and `server.key` in the root folder.
3. **Start Services**:
 ```bash
 docker-compose up -d --build
 ```
4. **Verify**:
 ```bash
 docker ps
 docker logs -f frostbyte-api-api-1
 ```

## 3. GitHub Actions (CI/CD)
To automate deployment, configure the following **Secrets** in your GitHub Repository settings:

- `DOCKER_USERNAME`: Your Docker Hub username.
- `DOCKER_PASSWORD`: Your Docker Hub token/password.
- `HOST`: The IP address of your production server.
- `USERNAME`: The SSH username (e.g., `ubuntu` or `root`).
- `SSH_KEY`: The private SSH key to access your server.

**Workflow Example (`.github/workflows/deploy.yml`):**
(You can add a workflow file to build and push the Docker image to Docker Hub, then SSH into your server to pull and restart.)

## 4. Security Checklist
- [ ] **HTTPS**: Ensure valid certificates are used.
- [ ] **Secrets**: Check `.env` is NOT in source control.
- [ ] **Firewall**: Allow only ports 80, 443, and SSH (22). Block 3306 (Database) from external access.
- [ ] **Rate Limiting**: Enabled in the application.
