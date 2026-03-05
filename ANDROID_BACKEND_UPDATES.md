# FrostByte API - Android Backend Updates

**Date**: January 21, 2026
**Status**: Completed
**Focus**: Stability, Performance, and Android App Integration

This document summarizes the critical backend updates implemented to resolve Android app crashes, data inconsistencies, and security issues.

---

## 1. New Endpoints

### Get Current User
* **Endpoint**: `GET /api/v1/user/me`
* **Purpose**: Allows the Android app to validate the session and fetch the current user's profile/roles without needing to store the User ID locally.
* **Response**: JSON object of the logged-in user.

### Dashboard Analytics
* **Endpoint**: `GET /api/v1/analytics/dashboard`
* **Purpose**: Provides pre-calculated stats for the dashboard to replace the client-side math that was broken by pagination.
* **Response Payload**:
 ```json
 {
 "total_revenue": 15420.50,
 "total_orders": 1250,
 "pending_count": 5,
 "daily_revenue": [
 {"date": "2026-01-21", "revenue": 450.00, "order_count": 12},
 ...
 ]
 }
 ```

---

## 2. API Enhancements

### Global Search
* **Endpoint**: `GET /api/v1/orders?search=...`
* **Functionality**: Fuzzy matches against **Order ID** or **Username**.
* **Why**: Enables finding specific orders even if they are not on the current page of results.

### Pagination & Filtering
* **Endpoint**: `GET /api/v1/orders`
* **New Parameters**:
 * `?page=1` (Default: 1)
 * `?limit=50` (Default: 50) - **Prevents app crashes from large payloads.**
 * `?status=PENDING,COOKING` (Filter by status)
 * `?date=YYYY-MM-DD` (Filter by date)

### Status Updates
* **Endpoint**: `PUT /api/v1/orders/{id}/status`
* **Improvement**: Now returns the **full updated Order object** instead of a generic success message.
* **Benefit**: Allows the Android app to update the specific item in the list immediately without refetching the entire list.

### Rate Limiting
* **Algorithm**: Token Bucket
* **Behavior**: Allows bursts of up to **10 requests/second** (smooths out "Click Spam" by cashiers) while maintaining an average rate limit.
* **Response**: Returns standardized JSON error (`429 Too Many Requests`) for better UI handling.

---

## 3. Authentication Upgrade

### Dual-Token Flow
To fix "Token valid: false" errors and improve security:
1. **Login (`POST /auth/login`)**:
 * Returns **two** tokens:
 * `access_token`: Short-lived (**15 minutes**). Used for API requests.
 * `refresh_token`: Long-lived (**7 days**). Used *only* to get new access tokens.
2. **Refresh (`POST /auth/refresh`)**:
 * **Payload**: `{ "refresh_token": "..." }` (Now in Body, not Header)
 * **Response**: `{ "access_token": "..." }`
 * **Logic**: Validates the refresh token signature and type before issuing a fresh access token.

---

## 4. Real-Time (WebSockets)

### Standardized Payload
* **Event**: `ORDER_UPDATED`
* **Payload Format**:
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
* **Trigger**: Sent to all Admin/Cashier devices whenever an order status changes.

---

## Code Cleanup
* Removed debug logging (route dumping) for cleaner startup.
* Standardized comments and removed TODOs.
* Ensured deployment readiness.

---

**Next Steps for Client**:
1. Update `Retrofit` or `OkHttp` interceptor to handle 401s by calling `/auth/refresh`.
2. Switch Dashboard to use `/analytics/dashboard`.
3. Implement Pagination support (`page++` on scroll).
