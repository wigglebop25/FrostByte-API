# Azure Deployment Guide for FrostByte API

## Prerequisites
- Azure CLI installed (`az login`)
- Docker installed

## Option 1: Azure Container Instances (Simple)

1. **Create a Resource Group**
 ```bash
 az group create --name FrostByteGroup --location eastus
 ```

2. **Create Azure Container Registry (ACR)**
 ```bash
 az acr create --resource-group FrostByteGroup --name frostbyteregistry --sku Basic --admin-enabled true
 ```

3. **Login to ACR**
 ```bash
 az acr login --name frostbyteregistry
 ```

4. **Build and Push Image**
 ```bash
 docker build -t frostbyteregistry.azurecr.io/frostbyte-api:v1 .
 docker push frostbyteregistry.azurecr.io/frostbyte-api:v1
 ```

5. **Deploy to Container Instance**
 ```bash
 az container create \
 --resource-group FrostByteGroup \
 --name frostbyte-api-container \
 --image frostbyteregistry.azurecr.io/frostbyte-api:v1 \
 --dns-name-label frostbyte-api \
 --ports 8080 \
 --environment-variables \
 DB_HOST=your-mysql-host \
 DB_USER=your-user \
 DB_PASSWORD=your-password \
 JWT_SECRET=your-production-secret
 ```

## Option 2: Azure App Service (Production)

1. Create a "Web App for Containers" in Azure Portal.
2. Select Docker Hub or Azure Container Registry.
3. Choose the image `frostbyteregistry.azurecr.io/frostbyte-api:v1`.
4. In **Configuration** > **Application Settings**, add your ENVs (`DB_HOST`, etc.).
5. Enable "Continuous Deployment" if using GitHub Actions.

## Database Note
For production, create an **Azure Database for MySQL** instance instead of running MySQL in a container side-by-side (though you can use Docker Compose in App Service, managed SQL is safer).