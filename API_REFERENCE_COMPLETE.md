# FrostByte API Reference

This document provides a comprehensive reference for all available endpoints in the FrostByte API v1.

**Base URL**: `https://frostbyte-api.southeastasia.cloudapp.azure.com/api/v1`

## Authentication & Security

* **Auth Type**: Bearer Token (JWT)
* **Header**: `Authorization: Bearer <your_token>`
* **Rate Limiting**: Token Bucket (allows bursts). 429 errors return structured JSON.
* **Security**: HTTPS/TLS enforced.

---

## 1. Authentication

### Get Current User (Me)
* **GET** `/user/me`
* **Access**: Authenticated Users
* **Response**: Returns the profile of the currently logged-in user (based on token).

### Register a new user
* **POST** `/auth/register`
* **Access**: Public
* **Payload**:
 ```json
 {
 "username": "jdoe",
 "password": "securePassword123"
 }
 ```

### Login
* **POST** `/auth/login`
* **Access**: Public
* **Payload**:
 ```json
 {
 "username": "jdoe",
 "password": "securePassword123"
 }
 ```
* **Response**: Returns `access_token` and `refresh_token`.

### Refresh Token
* **POST** `/auth/refresh`
* **Access**: Public (requires valid refresh token in body)
* **Payload**:
 ```json
 {
 "refresh_token": "eyJ..."
 }
 ```

---

## 2. Products

### Get All Products
* **GET** `/products`
* **Access**: Authenticated Users
* **Response**: List of all products with their categories.

### Get Product by ID
* **GET** `/products/{id}`
* **Access**: Authenticated Users

### Create Product
* **POST** `/products`
* **Access**: Admin Only
* **Payload**:
 ```json
 {
 "name": "Spicy Tuna Roll",
 "description": "Fresh tuna with spicy mayo",
 "price": 8.99,
 "product_image_uri": "/images/spicy-tuna.jpg"
 }
 ```

### Update Product
* **PUT** `/products/{id}`
* **Access**: Admin Only

### Delete Product
* **DELETE** `/products/{id}`
* **Access**: Admin Only

---

## 3. Categories

### Get All Categories
* **GET** `/categories`
* **Access**: Authenticated Users

### Get Category by ID
* **GET** `/categories/{id}`
* **Access**: Authenticated Users

### Create/Update/Delete Category
* **POST** `/categories` | **PUT** `/categories/{id}` | **DELETE** `/categories/{id}`
* **Access**: Admin Only

### Manage Category Products
* **POST** `/categories/product` (Add product to category)
* **POST** `/categories/product/remove` (Remove product from category)
* **Access**: Admin Only
* **Payload**:
 ```json
 {
 "category_id": 1,
 "product_id": 5
 }
 ```

---

## 4. Orders & Analytics

### Create Order
* **POST** `/orders`
* **Access**: Authenticated Users (Customers)
* **Payload**:
 ```json
 {
 "items": [
 { "product_id": 1, "quantity": 2 },
 { "product_id": 3, "quantity": 1 }
 ]
 }
 ```

### Get My Orders (or All Orders for Staff)
* **GET** `/orders`
* **Query Params**:
 * `?page=1` (Default: 1)
 * `?limit=50` (Default: 50)
 * `?status=PENDING,COOKING` (Filter by status)
 * `?date=2026-01-21` (Filter by date)
* **Access**: Authenticated Users
* **Behavior**:
 * **Customers**: Returns only their own orders.
 * **Admin/Cashier**: Returns ALL orders (system-wide), filtered by params.

### Get Order by ID
* **GET** `/orders/{id}`
* **Access**: Owner of order, Admin, or Cashier.

### Get User's Orders
* **GET** `/orders/user/{username}`
* **Access**: Admin, Cashier, or the user themselves.

### Update Order Status
* **PUT** `/orders/{id}/status`
* **Access**: Admin, Cashier
* **Payload**:
 ```json
 {
 "status": "COMPLETED" // PENDING, PREPARING, READY, COMPLETED, CANCELLED
 }
 ```
* **Response**: Returns the full updated `Order` object.

### Revenue Analytics (New!)
* **GET** `/analytics/revenue`
* **Access**: Admin, Cashier
* **Response**: Detailed revenue breakdown (Total Revenue, Total Orders, Average Order Value).

### Order Analytics (New!)
* **GET** `/orders/analytics`
* **Access**: Admin, Cashier
* **Response**: List of daily stats (Date, Total Orders, Completed Orders, Total Revenue).

---

## 5. Users & Roles

### Get All Users
* **GET** `/users`
* **Access**: Admin Only

### Search Users
* **GET** `/users/search?q=john`
* **Access**: Admin Only

### Manage Roles
* **GET** `/roles`
* **POST** `/roles/create`
* **POST** `/roles/assign` (Assign role to user)
* **POST** `/roles/add_permission`
* **Access**: Admin Only

---

## 6. Real-Time (WebSockets)

### Connect
* **URL**: `wss://frostbyte-api.southeastasia.cloudapp.azure.com/ws`
* **Auth**: Send JWT token in the query string `?token=YOUR_JWT_TOKEN` or as the first message.

### Events Received
The server pushes JSON events for live updates:
* `ORDER_UPDATED`: Status changed. Payload:
 ```json
 {
 "event": "ORDER_UPDATED",
 "data": {
 "order_id": 123,
 "status": "READY",
 "updated_at": "...",
 "order": { ...full_order_object... }
 }
 }
 ```

---

*Generated for FrostByte API v1 - Production*
