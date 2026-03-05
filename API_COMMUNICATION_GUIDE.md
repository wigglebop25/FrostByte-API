# FrostByte API Communication Guide

This guide details how the Android client should communicate with the FrostByte API.

## Base URLs
* **Production:** `https://frostbyte-api.southeastasia.cloudapp.azure.com/api/v1`
* **WebSocket:** `wss://frostbyte-api.southeastasia.cloudapp.azure.com/ws`

## Authentication
**All endpoints (except Login/Register) require a JWT Token.**

### 1. Login
* **Endpoint:** `POST /auth/login`
* **Body:**
 ```json
 {
 "username": "customer",
 "password": "customer123"
 }
 ```
* **Response:**
 ```json
 {
 "token": "eyJhbGciOiJIUzI1NiIsIn...",
 "user_id": 3,
 "role": "Customer"
 }
 ```
* **Action:** Save the `token` securely (e.g., EncryptedSharedPreferences).

## Making Requests
Add the token to the header of **every** HTTP request.

* **Header Name:** `Authorization`
* **Header Value:** `Bearer <YOUR_SAVED_TOKEN>`

### Example: Get Products (Menu)
* **Endpoint:** `GET /products`
* **Response:**
 ```json
 [
 {
 "product_id": 1,
 "name": "Pork Gyoza",
 "price": 8.50,
 "product_image_uri": "/images/pork-gyoza.jpg",
 "description": "Six pieces of pan-seared dumplings..."
 }
 ]
 ```

### Example: Place Order
* **Endpoint:** `POST /orders`
* **Body:**
 ```json
 {
 "products": [
 { "product_id": 1, "quantity": 2 },
 { "product_id": 4, "quantity": 1 }
 ]
 }
 ```

### Example: Get My Orders (History)
* **Endpoint:** `GET /orders`
* **Response:**
 ```json
 [
 {
 "order_id": 105,
 "status": "PENDING",
 "total_amount": 29.00,
 "created_at": "January 17, 2026 04:32 PM",
 "products": [...]
 }
 ]
 ```

## WebSocket (Real-time Updates)
Use this to listen for order status changes (e.g., when an order becomes `READY`).

* **Connection URL:**
 `wss://frostbyte-api.southeastasia.cloudapp.azure.com/ws?token=<YOUR_SAVED_TOKEN>`
 *(Note: Token is passed in the URL query, NOT the header)*

* **Events Received:**
 ```json
 {
 "order_id": 105,
 "status": "READY",
 "updated_at": "January 17, 2026 04:45 PM"
 }
 ```

## Timezone Handling
* The API returns all times in **Philippine Time (UTC+8)**.
* Format: `"January 17, 2026 04:32 PM"`
* **Android Action:** You can display this string directly to the user. No timezone conversion is needed on the phone.

## Common Errors
* **401 Unauthorized:** Your token is missing, expired, or invalid. -> *Redirect user to Login.*
* **403 Forbidden:** You are trying to access an Admin resource (like `/users`) as a Customer.
* **404 Not Found:** You might be using the wrong ID or URL.
