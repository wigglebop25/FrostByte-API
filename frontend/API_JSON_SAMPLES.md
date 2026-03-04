# FrostByte API - JSON Samples

This document provides ready-to-use JSON samples for every API endpoint to help with frontend integration and mocking.

---

## 1. Authentication

### **POST** `/auth/register`
**Request:**
```json
{
 "username": "cool_cat_99",
 "password": "SecurePassword123!"
}
```
**Response (201 Created):**
```json
{
 "message": "Account created successfully"
}
```

### **POST** `/auth/login`
**Request:**
```json
{
 "username": "cool_cat_99",
 "password": "SecurePassword123!"
}
```
**Response (200 OK):**
```json
{
 "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJyb2xlcyI6WyJDdXN0b21lciJdfQ.sig",
 "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0eXBlIjoicmVmcmVzaCJ9.sig",
 "user": {
 "user_id": 1,
 "username": "cool_cat_99",
 "roles": [
 {
 "role_id": 2,
 "name": "Customer",
 "permissions": "READ",
 "description": "Standard user",
 "created_at": "January 25, 2026 10:00 AM",
 "updated_at": "January 25, 2026 10:00 AM"
 }
 ],
 "created_at": "January 25, 2026 10:00 AM",
 "updated_at": "January 25, 2026 10:00 AM"
 }
}
```

### **POST** `/auth/refresh`
**Request:**
```json
{
 "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0eXBlIjoicmVmcmVzaCJ9.sig"
}
```
**Response (200 OK):**
```json
{
 "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.new_access_token.sig"
}
```

---

## 2. User Management

### **GET** `/user/me`
**Response (200 OK):**
```json
{
 "user_id": 1,
 "username": "cool_cat_99",
 "roles": [
 {
 "role_id": 2,
 "name": "Customer"
 }
 ],
 "created_at": "January 25, 2026 10:00 AM",
 "updated_at": "January 25, 2026 10:00 AM"
}
```

---

## 3. Products & Categories

### **GET** `/products`
**Response (200 OK):**
```json
[
 {
 "product_id": 101,
 "name": "California Roll",
 "product_image_uri": "/products/cali_roll.png",
 "description": "Crab meat, avocado, and cucumber wrapped in seaweed and rice.",
 "price": 8.50,
 "categories": [
 {
 "category_id": 1,
 "name": "Sushi",
 "description": "Traditional Japanese rice dishes"
 }
 ],
 "created_at": "January 20, 2026 09:30 AM",
 "updated_at": "January 20, 2026 09:30 AM"
 },
 {
 "product_id": 102,
 "name": "Tonkotsu Ramen",
 "product_image_uri": "/products/tonkotsu_ramen.png",
 "description": "Rich pork broth with noodles, chashu, and egg.",
 "price": 12.99,
 "categories": [
 {
 "category_id": 2,
 "name": "Ramen",
 "description": "Noodle soups"
 }
 ],
 "created_at": "January 20, 2026 09:35 AM",
 "updated_at": "January 20, 2026 09:35 AM"
 }
]
```

### **POST** `/products` (Admin Only)
**Request:**
```json
{
 "name": "Spicy Salmon Roll",
 "description": "Fresh salmon with spicy mayo and cucumber",
 "price": 9.50,
 "product_image_uri": "/products/spicy_salmon.png",
 "categories": ["Sushi", "Spicy"]
}
```

### **GET** `/categories`
**Response (200 OK):**
```json
[
 {
 "category_id": 1,
 "name": "Sushi",
 "description": "Traditional Japanese rice dishes",
 "products": null, 
 "created_at": "January 15, 2026 08:00 AM",
 "updated_at": "January 15, 2026 08:00 AM"
 },
 {
 "category_id": 2,
 "name": "Beverages",
 "description": "Cold and hot drinks",
 "products": null,
 "created_at": "January 15, 2026 08:05 AM",
 "updated_at": "January 15, 2026 08:05 AM"
 }
]
```

---

## 4. Orders

### **POST** `/orders` (Create Order)
**Request:**
```json
{
 "products": [
 {
 "product_id": 101,
 "quantity": 2
 },
 {
 "product_id": 102,
 "quantity": 1
 }
 ]
}
```
**Response (201 Created):**
```json
{
 "order_id": 505,
 "user_id": 1,
 "total_amount": 29.99,
 "status": "PENDING",
 "products": [
 {
 "order_id": 505,
 "product_id": 101,
 "quantity": 2,
 "unit_price": 8.50,
 "line_total": 17.00
 },
 {
 "order_id": 505,
 "product_id": 102,
 "quantity": 1,
 "unit_price": 12.99,
 "line_total": 12.99
 }
 ],
 "created_at": "January 25, 2026 12:45 PM",
 "updated_at": "January 25, 2026 12:45 PM"
}
```

### **GET** `/orders`
**Response (200 OK):**
```json
[
 {
 "order_id": 505,
 "user_id": 1,
 "total_amount": 29.99,
 "status": "PENDING",
 "products": [
 {
 "product_id": 101,
 "product": {
 "name": "California Roll",
 "product_image_uri": "/products/cali_roll.png",
 "price": 8.50
 },
 "quantity": 2,
 "unit_price": 8.50,
 "line_total": 17.00
 }
 ],
 "created_at": "January 25, 2026 12:45 PM",
 "updated_at": "January 25, 2026 12:45 PM"
 }
]
```

### **PUT** `/orders/{id}/status` (Staff Only)
**Request:**
```json
{
 "status": "READY"
}
```
**Response (200 OK):**
```json
{
 "order_id": 505,
 "status": "READY",
 "updated_at": "January 25, 2026 01:00 PM"
 // ... other order fields
}
```

---

## 5. Analytics

### **GET** `/analytics/dashboard`
**Response (200 OK):**
```json
{
 "total_revenue": 15430.50,
 "total_orders": 342,
 "pending_orders": 5,
 "completed_orders": 330,
 "cancelled_orders": 7,
 "average_order_value": 45.12,
 "daily_revenue": {
 "2026-01-20": 2100.50,
 "2026-01-21": 1850.00,
 "2026-01-22": 2300.75,
 "2026-01-23": 2500.00,
 "2026-01-24": 3100.25,
 "2026-01-25": 1200.00
 }
}
```

---

## 6. WebSockets

### **Event: New Order** (Received by Admin/Cashier)
*Note: This is the raw Order object.*
```json
{
 "order_id": 506,
 "user_id": 42,
 "total_amount": 45.00,
 "status": "PENDING",
 "products": [
 {
 "product_id": 105,
 "quantity": 3,
 "unit_price": 15.00,
 "line_total": 45.00
 }
 ],
 "created_at": "January 25, 2026 01:15 PM",
 "updated_at": "January 25, 2026 01:15 PM"
}
```

### **Event: Order Updated** (Received by Customer & Staff)
```json
{
 "event": "ORDER_UPDATED",
 "data": {
 "order_id": 505,
 "status": "READY",
 "updated_at": "January 25, 2026 01:20 PM",
 "order": {
 "order_id": 505,
 "user_id": 1,
 "status": "READY",
 "total_amount": 29.99
 // ... full order object
 }
 }
}
```
