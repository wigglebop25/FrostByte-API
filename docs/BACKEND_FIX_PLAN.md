# Backend Fix Plan

## Problem Analysis

The Android app is experiencing `401 Unauthorized` and `403 Forbidden` errors.

### 1. 401 Unauthorized
*   **Cause:** Token invalidation due to server resets or `JWT_SECRET` changes.
*   **Fix:** Ensure the Android app logs out and logs back in to get a fresh token signed with the current secret.

### 2. 403 Forbidden (Orders)
*   **Endpoint:** `GET /api/v1/orders/user/{username}`
*   **Root Cause:** The route was incorrectly nested inside a `RoleMiddleware("Admin", "Cashier")` block in `cmd/api/main.go`. This prevented Customers from accessing their own orders, even though the Handler logic was updated to allow it.
*   **Fix:** Move the route out of the middleware group and rely on the Handler's internal security check (Self-Access or Staff).

### 3. 403 Forbidden (Users)
*   **Endpoint:** `GET /api/v1/users/{id}`
*   **Root Cause:** The entire `/users` group was protected by `AdminMiddleware`.
*   **Fix:** 
    1.  Move `GET /{id}` out of the Admin-only group.
    2.  Update `UserHandler.GetByID` to implement self-access security (allow if `token.UserID == requestedUserID`).

## Implementation Steps

1.  **Modify `cmd/api/main.go`**:
    *   Refactor `/orders` routes to expose `GetByUser` to authenticated users (logic inside handler protects it).
    *   Refactor `/users` routes to expose `GetByID` to authenticated users.

2.  **Modify `internal/handlers/user_handler.go`**:
    *   Update `GetByID` to check if the requester is the user themselves or an Admin.

3.  **Verify `internal/handlers/order_handler.go`**:
    *   Confirm `GetByUser` has the logic to check `username == tokenUsername`. (Already done).

## Token Claims
*   Added `sub` claim to JWT to support standard client libraries.
*   Added `username` and `roles` to JWT.

This plan ensures that while endpoints become "accessible" at the router level, the data remains strictly protected by business logic in the handlers.
