# Master Migration Plan: FrostByte API

## 🎯 Objectives
- **Language:** Go (Golang) 1.23+
- **Architecture:** Clean Architecture / Handler-Service-Repository pattern.
- **Database:** MySQL (preserving existing schema structure).
- **ORM:** GORM (for rapid development) or SQLx (for performance). *Decision: GORM for this plan to match velocity.*
- **Real-time:** WebSockets for live transaction/order updates.
- **Testing:** Native Go tests + Bruno for API integration tests.
- **Deployment:** Docker + Azure.
- **Security:** Standard HTTPS, JWT Authentication, Role-Based Access Control (RBAC).

## 🛠 Tech Stack
- **Framework:** [Chi](https://github.com/go-chi/chi) (Lightweight, idiomatic).
- **Database:** MySQL.
- **ORM:** [GORM](https://gorm.io/).
- **WebSockets:** [Gorilla WebSocket](https://github.com/gorilla/websocket).
- **Config:** [Godotenv](https://github.com/joho/godotenv).

---

## 📅 Implementation Phases

### Phase 1: Foundation & Infrastructure 🏗️
**Goal:** Initialize the Go project, set up Docker, and establish database connectivity.

- [x] **1.1. Project Initialization:**
    - `go mod init frostbyte-api`
    - Setup folder structure (`cmd`, `internal`, `pkg`, `configs`, `api`).
- [x] **1.2. Docker Setup:**
    - Create `compose.yml` (MySQL + Go App).
    - Create `Dockerfile` for the Go application (Multi-stage build).
- [x] **1.3. Configuration Management:**
    - Implement environment variable loading (database credentials, JWT secrets, port).
- [x] **1.4. Database Connection:**
    - Initialize GORM connection to MySQL.

### Phase 2: Domain Modeling & Auth Security 🔐
**Goal:** Replicate User and Role management with secure Authentication.

- [x] **2.1. Domain Models:**
    - Define `User`, `Role`, `UserRole` structs in Go (matching `schema.rs`).
- [x] **2.2. Migrations:**
    - Setup GORM auto-migration to replicate the schema.
- [x] **2.3. Repository Layer (Users/Roles):**
    - Implement methods for CRUD on Users and Roles.
- [x] **2.4. Authentication Service:**
    - Implement Password Hashing (Argon2 or bcrypt).
    - Implement JWT generation and validation (Access & Refresh tokens).
- [x] **2.5. HTTP Handlers (Auth):**
    - `/auth/login`, `/auth/register`, `/auth/refresh`.
- [x] **2.6. Middleware:**
    - JWT Authentication Middleware (Basic setup in Handler, expandable).

### Phase 3: Inventory Management (Products & Categories) 📦
**Goal:** CRUD for Products and Categories with relationships.

- [x] **3.1. Domain Models:**
    - Define `Category`, `Product`, `ProductCategory`.
- [x] **3.2. Repository Layer:**
    - Implement CRUD for Products and Categories.
    - Handle Many-to-Many relationships.
- [x] **3.3. HTTP Handlers:**
    - `/products` (GET, POST, PUT, DELETE).
    - `/categories` (GET, POST, PUT, DELETE).
    - Manage product-category assignments.

### Phase 4: Orders & Transactions 🛒
**Goal:** Handle Order creation and management.

- [x] **4.1. Domain Models:**
    - Define `Order`, `OrderProduct`.
- [x] **4.2. Repository Layer:**
    - Implement Transactional Order Creation (ensure `orders` and `order_products` are saved atomically).
- [x] **4.3. Service Layer:**
    - Calculate totals.
    - Validate stock/availability (if applicable).
- [x] **4.4. HTTP Handlers:**
    - `/orders` (Create, Get All, Get By ID).
    - Update Order Status.

### Phase 5: Real-time Live Updates ⚡
**Goal:** Broadcast new orders/status changes to connected clients (e.g., Kitchen Display).

- [x] **5.1. WebSocket Hub:**
    - Implement a WebSocket Hub to manage connections (Clients).
- [x] **5.2. Integration:**
    - Hook into the `OrderService`. When an order is created or status changes -> Broadcast message.
- [x] **5.3. WebSocket Endpoint:**
    - `/ws/orders` endpoint.

### Phase 6: Testing & Quality Assurance 🧪
**Goal:** Verify system stability and correctness.

- [x] **6.1. Unit Tests:**
    - Test Services and Utility functions (`internal/service/order_service_test.go`).
- [x] **6.2. Bruno API Tests:**
    - Created `bruno/go-api` with Health and Login requests.
- [x] **6.3. Security Audit:**
    - Ensure HTTPS headers.
    - Validate input sanitization (Basic JSON decoding validation).

### Phase 7: Deployment & Polishing 🚀
**Goal:** Prepare for Azure deployment.

- [x] **7.1. Final Docker Optimization:**
    - Multi-stage build implemented in `Dockerfile`.
- [x] **7.2. Azure Configuration:**
    - Updated `AZURE_DEPLOY.md` for **FrostByte API**.
- [x] **7.3. Documentation:**
    - `README.md` updated for **FrostByte API**.

---

## 🚦 Tracking Progress

All phases completed. The server is ready to run via `docker-compose up --build`.

---

# 📚 Appendix: Implementation Reference

Use these details to reconstruct the code exactly as it was.

## A. Dependencies (`go.mod`)

```bash
go get github.com/go-chi/chi/v5
go get github.com/go-chi/cors
go get github.com/golang-jwt/jwt/v5
go get github.com/joho/godotenv
go get golang.org/x/crypto
go get gorm.io/driver/mysql
go get gorm.io/gorm
go get github.com/gorilla/websocket
```

## B. Environment Variables (`.env`)

| Variable | Default | Description |
| :--- | :--- | :--- |
| `SERVER_PORT` | `8080` | Port for the API server |
| `DB_HOST` | `localhost` | Database hostname |
| `DB_PORT` | `3306` | Database port |
| `DB_USER` | `root` | Database username |
| `DB_PASSWORD` | (empty) | Database password |
| `DB_NAME` | `frost_db` | Name of the database schema |
| `JWT_SECRET` | `default_secret` | Secret key for JWT |

## C. Domain Models (`internal/domain/models.go`)

```go
type User struct {
    UserID       uint      `gorm:"primaryKey;column:user_id" json:"user_id"`
    Username     string    `gorm:"size:50;not null;unique" json:"username"`
    PasswordHash string    `gorm:"size:255;not null" json:"-"`
    Roles        []Role    `gorm:"many2many:user_roles;joinForeignKey:user_id;joinReferences:role_id" json:"roles,omitempty"`
    CreatedAt    time.Time `json:"created_at"`
    UpdatedAt    time.Time `json:"updated_at"`
}

type Product struct {
    ProductID       uint       `gorm:"primaryKey;column:product_id" json:"product_id"`
    Name            string     `gorm:"size:100;not null" json:"name"`
    Price           float64    `gorm:"type:decimal(10,2);not null" json:"price"`
    Categories      []Category `gorm:"many2many:product_categories;joinForeignKey:product_id;joinReferences:category_id" json:"categories,omitempty"`
}

type Order struct {
    OrderID     uint           `gorm:"primaryKey;column:order_id" json:"order_id"`
    UserID      uint           `gorm:"column:user_id;not null" json:"user_id"`
    TotalAmount float64        `gorm:"type:decimal(10,2);not null" json:"total_amount"`
    Status      string         `gorm:"size:50;default:'PENDING'" json:"status"`
    Products    []OrderProduct `gorm:"foreignKey:OrderID" json:"products,omitempty"`
}
// ... (Include other structs as seen in Appendix C previously)
```

## D. WebSocket Implementation Details (`internal/websocket/hub.go`)

To implement real-time updates, use a **Hub** to manage connections.

1.  **Hub Struct**: Manages `clients` (map), `broadcast` (chan), `register` (chan), and `unregister` (chan).
2.  **Run() Method**: A loop in a goroutine that selects from these channels.
3.  **Broadcast() Method**: Sends a message to the `broadcast` channel.
4.  **Integration**: Initialize in `main.go`, pass to `OrderService`, and call `Broadcast` when an order is created.

## E. Database Initialization Strategy

**Method 1: GORM Auto-Migration (Recommended)**
In `main.go`, call `db.AutoMigrate(&User{}, &Role{}, &Product{}, &Category{}, &Order{}, &OrderProduct{})`.

**Method 2: Manual SQL**
Refer to the included **`guidance/DATABASE_SCHEMA.sql`** for the full DDL.