# FrostByte API Documentation (Web Frontend)

**Version:** v1
**Base URL:** `https://frostbyte-api.southeastasia.cloudapp.azure.com/api/v1`

---

## CRITICAL FRONTEND NOTES

1. **Authentication:**
 * Use **USERNAME** for login/register. **Do NOT use Email.**
 * Store the received `token` (JWT) in `localStorage` or `sessionStorage`.
 * Include the token in the header of **every** authenticated request:
 `Authorization: Bearer <your_token>`

2. **CORS & Environment:**
 * The API accepts requests from any origin (`*`).
 * Ensure your `VITE_API_URL` (or equivalent) points to the Base URL above.

3. **Date Formatting:**
 * Timestamps are returned as strings: `"January 02, 2006 03:04 PM"`.
 * For date filtering (query params), use `YYYY-MM-DD`.

---

## 1. Authentication

### Register
* **POST** `/auth/register`
* **Payload**:
 ```json
 {
 "username": "myuser", // REQUIRED
 "password": "mypassword" // REQUIRED
 }
 ```

### Login
* **POST** `/auth/login`
* **Payload**:
 ```json
 {
 "username": "myuser",
 "password": "mypassword"
 }
 ```
* **Response**:
 ```json
 {
 "token": "eyJ...", // Save this!
 "refresh_token": "eyJ...", // Save this for refreshing
 "user": { ... }
 }
 ```

### Get Current User
* **GET** `/user/me`
* **Headers**: `Authorization: Bearer <token>`
* **Response**: `User` object.

---

## 2. Products

### Get All Products
* **GET** `/products`
* **Response**: `Product[]`
 ```json
 [
 {
 "product_id": 1,
 "name": "California Roll",
 "price": 8.50,
 "product_image_uri": "/products/cali_roll.png",
 "categories": [ { "name": "Sushi" } ]
 }
 ]
 ```

### Create Product (Admin)
* **POST** `/products`
* **Payload**:
 ```json
 {
 "name": "Spicy Tuna",
 "description": "Spicy tuna roll",
 "price": 9.00,
 "product_image_uri": "/products/spicy_tuna.png",
 "categories": ["Sushi", "Spicy"] // Array of Category NAMES
 }
 ```

---

## 3. Orders

### Create Order (Checkout)
* **POST** `/orders`
* **Headers**: `Authorization: Bearer <token>`
* **Payload**:
 ```json
 {
 "products": [ // KEY IS "products", NOT "items"
 {
 "product_id": 1,
 "quantity": 2
 },
 {
 "product_id": 5,
 "quantity": 1
 }
 ]
 }
 ```

### Get Orders (List)
* **GET** `/orders`
* **Query Params**:
 * `?page=1` & `?limit=50`
 * `?status=PENDING` (or `READY`, `COMPLETED`, `CANCELLED`)
 * `?date=2026-01-25`
* **Response**: `Order[]`

### Update Status (Admin/Cashier)
* **PUT** `/orders/{id}/status`
* **Payload**:
 ```json
 {
 "status": "READY"
 }
 ```

---

## 4. Analytics (Admin/Cashier)

### Dashboard Stats
* **GET** `/analytics/dashboard`
* **Response**:
 ```json
 {
 "total_revenue": 5000.50,
 "total_orders": 120,
 "pending_orders": 5,
 "completed_orders": 110,
 "average_order_value": 45.45,
 "daily_revenue": {
 "2026-01-24": 1200.00,
 "2026-01-25": 1500.50
 }
 }
 ```

---

## 5. WebSockets (Real-Time)

### Connection
* **URL**: `wss://frostbyte-api.southeastasia.cloudapp.azure.com/ws`
* **Auth**: Append token as query param.
 * `wss://...?token=<YOUR_JWT_TOKEN>`

### Client Implementation Guide (Svelte/React)
1. **Connect**: Open socket with token.
2. **Ping**: Send `{"type": "ping"}` every 30s to keep connection alive.
3. **Listen**: Handle incoming JSON messages.

### Incoming Events
The server sends a JSON object with an `event` type and `data`.

**1. Order Updated** (Status change)
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

**2. New Order** (For Staff Dashboard)
```json
{
 "event": "NEW_ORDER",
 "data": {
 "order_id": 124,
 "total_amount": 50.00,
 "status": "PENDING",
 ...
 }
}
```
