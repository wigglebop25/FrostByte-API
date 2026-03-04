# FrostByte API - Frontend Integration Guide

This document is generated from the **Bruno** collection and represents the complete, actionable API surface for the Frontend. It covers all workflows including Authentication, Admin Management (modifying data), Cashier operations, and Customer actions.

**Base URL**: `https://frostbyte-api.southeastasia.cloudapp.azure.com/api/v1`
**Auth**: Most endpoints require a **Bearer Token**. Include `Authorization: Bearer <token>` in headers.

---

## 0. Authentication
*These endpoints handle user entry. Store the returned `token` for subsequent requests.*

### Login
* **POST** `/auth/login`
* **Auth**: None
* **Payload**:
 ```json
 {
 "username": "adminuser",
 "password": "adminpassword"
 }
 ```
* **Response**: Returns `{ "token": "...", "user": {...} }`.

### Register
* **POST** `/auth/register`
* **Auth**: None
* **Payload**:
 ```json
 {
 "username": "newuser",
 "password": "securepassword"
 }
 ```

### Refresh Token
* **POST** `/auth/refresh`
* **Auth**: Bearer Token
* **Payload**:
 ```json
 {
 "refresh_token": "..." 
 }
 ```
 *(Note: Ensure you send the refresh token string)*

---

## 1. Admin Workflow (Management)
*Operations for "Modifying Everything" and "Adding Users". Requires Admin Role.*

### Users (Add & Modify)

#### **Add User**
Create a new user manually with a specific role.
* **POST** `/users`
* **Payload**:
 ```json
 {
 "username": "cashier1",
 "password": "password123",
 "role_name": "Cashier" // e.g., 'Admin', 'Cashier', 'Customer'
 }
 ```

#### **Edit User**
Modify an existing user's credentials or role.
* **PUT** `/users/{id}`
* **Payload**:
 ```json
 {
 "username": "updateduser",
 "password": "updatedpassword", // Send only if changing
 "role_name": "Cashier"
 }
 ```

#### **Get All Users**
* **GET** `/users`

#### **Get User By ID**
* **GET** `/users/{id}`

#### **Search User**
* **GET** `/users/search?username={name}`

#### **Delete User**
* **DELETE** `/users/{id}`

---

### Categories

#### **Create Category**
* **POST** `/categories`
* **Payload**:
 ```json
 {
 "name": "Appetizers",
 "description": "Starters"
 }
 ```

#### **Edit Category**
* **PUT** `/categories/{id}`
* **Payload**:
 ```json
 {
 "name": "Updated Name",
 "description": "Updated Description"
 }
 ```

#### **Delete Category**
* **DELETE** `/categories/{id}`

#### **Get All Categories**
* **GET** `/categories`

#### **Get Products by Category**
* **GET** `/categories/{category_name}/products`
* *Example*: `/categories/Main%20Course/products`

#### **Remove Product from Category**
* **POST** `/categories/product/remove`
* **Payload**:
 ```json
 {
 "category": "Food", // Category Name
 "product": "Burger" // Product Name
 }
 ```

---

### Products

#### **Create Product**
* **POST** `/products`
* **Payload**:
 ```json
 {
 "name": "Burger",
 "description": "Delicious beef burger",
 "price": 9.99,
 "product_image_uri": "/images/burger.jpg",
 "categories": ["Main Course", "Food"] // Array of Category Names
 }
 ```

#### **Update Product**
* **PUT** `/products/{id}`
* **Payload**:
 ```json
 {
 "name": "Double Burger",
 "description": "Double beef",
 "price": 14.99,
 "product_image_uri": "/images/double.jpg",
 "categories": ["Food"]
 }
 ```

#### **Delete Product**
* **DELETE** `/products/{id}`

#### **Get All Products**
* **GET** `/products`

---

### Roles & Permissions

#### **Create Role**
* **POST** `/roles/create`
* **Payload**:
 ```json
 {
 "name": "SuperAdmin",
 "description": "God mode"
 }
 ```

#### **Update Role**
* **POST** `/roles/update/{id}`
* **Payload**:
 ```json
 {
 "name": "superadmin",
 "description": "Updated description"
 }
 ```

#### **Assign Role to User**
* **POST** `/roles/assign`
* **Payload**:
 ```json
 {
 "username": "jdoe",
 "role_name": "Admin"
 }
 ```

#### **Add Permission to Role**
* **POST** `/roles/add_permission`
* **Payload**:
 ```json
 {
 "role_name": "customer",
 "permission": "WRITE"
 }
 ```

#### **Set Permission (Override)**
* **POST** `/roles/{id}/set_permission`
* **Payload**: `{ "permission": "READ" }`

#### **Remove Permission**
* **PATCH** `/roles/{id}/delete_permission`
* **Payload**: `{ "permission": "read" }`

#### **Delete Role**
* **DELETE** `/roles/{id}`

#### **Get All Roles**
* **GET** `/roles`

---

### Orders (Admin View)

#### **Get All Orders**
* **GET** `/orders`

#### **Get Order By ID**
* **GET** `/orders/{id}`

#### **Get Orders By User**
* **GET** `/orders/user/{username}`

#### **Get Orders By Role**
* **GET** `/orders/role/{role_name}`

#### **Get Analytics**
* **GET** `/orders/analytics`
* *Returns revenue stats, status counts, etc.*

---

## 2. Cashier Workflow

#### **View Queue (All Orders)**
* **GET** `/orders`

#### **Filter Orders by Customer**
* **GET** `/orders/user/{username}`

#### **Update Order Status**
* **PUT** `/orders/{id}/status`
* **Payload**:
 ```json
 {
 "status": "ready" // pending, preparing, ready, completed, cancelled
 }
 ```

#### **View Analytics**
* **GET** `/orders/analytics`

---

## 3. Customer Workflow

#### **View Menu (Products)**
* **GET** `/products`

#### **View Product Details**
* **GET** `/products/{id}`

#### **Place Order**
* **POST** `/orders`
* **Payload**:
 ```json
 {
 "products": [
 {
 "product_id": 1,
 "quantity": 2
 },
 {
 "product_id": 2,
 "quantity": 1
 }
 ]
 }
 ```

#### **View My History**
* **GET** `/orders`

---

## General

#### **Health Check**
* **GET** `/`
* *Returns 200 OK if API is running.*
