# FrostByte API: Implementation Walkthrough

This guide is your daily "instruction manual" for building the FrostByte API. It breaks down the implementation into specific coding tasks.

---

## 🏁 Session 1: The Skeleton (Setup)
**Goal:** Get a basic Go server running.

### 1. Initialize Project
Open your terminal in `D:\Projects\School\FrostByte API`:
```powershell
go mod init frostbyte-api
```

### 2. Install Dependencies
Copy and run this block:
```powershell
go get github.com/go-chi/chi/v5
go get github.com/go-chi/cors
go get github.com/golang-jwt/jwt/v5
go get github.com/joho/godotenv
go get golang.org/x/crypto
go get gorm.io/driver/mysql
go get gorm.io/gorm
go get github.com/gorilla/websocket
```

### 3. Create Folder Structure
Create these folders:
```text
cmd/api/
internal/config/
internal/database/
internal/domain/
internal/handlers/
internal/repository/
internal/service/
internal/websocket/
```

### 4. Create Configuration
**File:** `internal/config/config.go`
**Action:** Write the `Config` struct and `LoadConfig()` function to read environment variables (e.g., SERVER_PORT, DB_DSN, JWT_SECRET).

### 5. Create Entry Point
**File:** `cmd/api/main.go`
*Action:* Write a basic `main()` function that:
1. Loads Config.
2. Initializes a `chi.NewRouter()`.
3. Adds a simple "Hello World" route: `r.Get("/", ...)`
4. Starts the server: `http.ListenAndServe(...)`

**✅ Checkpoint:** Run `go run cmd/api/main.go`. Open `http://localhost:8080`. You should see your message.

---

## 🗄️ Session 2: The Heart (Database)
**Goal:** Connect to MySQL using Docker.

### 1. Docker Setup
**File:** `compose.yml` (in root)
**Action:** Define a MySQL service (port 3306) and the Go app service.

### 2. Domain Models
**File:** `internal/domain/models.go`
**Action:** Create structs for `User`, `Product`, `Order`, `Category`, `Role`, etc.

### 3. Database Connection
**File:** `internal/database/database.go`
*Action:* Implement `Connect(dsn string)`.
*   Use `gorm.Open` with the MySQL driver.
*   **Critical:** Add the Auto-Migration logic:
    ```go
    DB.AutoMigrate(&domain.User{}, &domain.Role{}, &domain.Product{}, &domain.Category{}, &domain.Order{}, &domain.OrderProduct{})
    ```

### 4. Wire it Up
**File:** `cmd/api/main.go`
*Action:* Add `database.Connect(cfg.DSN())` before starting the server.

**✅ Checkpoint:** Run `docker-compose up --build`. Check logs for "Database connected successfully!".

---

## 🔐 Session 3: The Gatekeeper (Auth)
**Goal:** Register users and issue JWTs.

### 1. Repository
**File:** `internal/repository/user_repo.go`
*Action:* Implement `Create(user)` and `FindByUsername(username)`.

### 2. Service
**File:** `internal/service/auth_service.go`
*Action:*
*   `Register`: Hash password (bcrypt), save user.
*   `Login`: Check password, generate JWT (golang-jwt).

### 3. Handler
**File:** `internal/handlers/auth_handler.go`
*Action:* Parse JSON request -> Call Service -> Return JSON response (token).

### 4. Routes
**File:** `cmd/api/main.go`
*Action:* Register routes `/api/v1/auth/login` and `/api/v1/auth/register`.

**✅ Checkpoint:** Open **Bruno**. Run "Auth > Register". You should get a 201 Created and a Token.

---

## 📦 Session 4: The Core (Business Logic)
**Goal:** Manage Products and Orders.

### 1. Products & Categories
*   **Repos:** `product_repo.go`, `category_repo.go` (CRUD operations).
*   **Services:** `product_service.go`, `category_service.go` (Pass-through logic).
*   **Handlers:** `product_handler.go`, `category_handler.go`.
*   **Routes:** Register `/products` and `/categories` in `main.go`.

### 2. Orders
*   **Repo:** `order_repo.go` (Create order + Save OrderProducts transactionally).
*   **Service:** `order_service.go` (Calculate totals, call Repo).
*   **Handler:** `order_handler.go` (Parse input, call Service).

**✅ Checkpoint:** Use Bruno to Create a Category, then a Product, then an Order. Check your Database to see the data linked correctly.

---

## ⚡ Session 5: The Pulse (WebSockets)
**Goal:** Live updates when orders are created.

### 1. The Hub
**File:** `internal/websocket/hub.go`
**Action:** Implement the WebSocket `Hub` to manage connections and broadcast messages.

### 2. Integration
**File:** `cmd/api/main.go`
*Action:* Initialize `hub := websocket.NewHub()`, run it `go hub.Run()`, and pass it to `NewOrderService`.

**File:** `internal/service/order_service.go`
*Action:* Update `CreateOrder` to call `s.hub.Broadcast(...)` after saving an order.

### 3. WebSocket Route
**File:** `internal/handlers/ws_handler.go`
*Action:* Add a handler to upgrade HTTP to WS. Register `/ws/orders` in `main.go`.

**✅ Checkpoint:** Open a WS connection in a tool (or browser console). Create an Order via Bruno. You should see the JSON message appear instantly in the WS connection.

---

## 🚀 Final Steps
1.  **Review `AZURE_DEPLOY.md`** when you are ready to go live.
2.  **Run `go test ./...`** to verify integrity.
