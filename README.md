# VELO Backend

REST API backend for **VELO**, an e-commerce application, written in **Go** (standard library `net/http`) following a _clean architecture_ pattern (Handler → Service → Repository).

Key features:

- 🔐 **Authentication**: registration + email verification, JWT login (stored in a cookie), logout, and **Google OAuth 2.0**
- 🛡️ **Authorization**: JWT and RBAC middleware (`admin` role required to manage products)
- 🛍️ **Products**: product CRUD with pagination and Redis caching (5-minute TTL)
- 🛒 **Cart**: add, view, update quantity, remove items, clear cart
- 💳 **Checkout & Payment**: **Midtrans** (Snap) integration with an **idempotency key** (stored in Redis) and a notification webhook
- ⏱️ **Scheduler**: automatically cancels expired `Unpaid` orders every 5 minutes
- 📧 **Email**: verification emails sent via **Resend**

---

## 🧰 Tech Stack

| Category       | Technology                                             |
| -------------- | ------------------------------------------------------ |
| Language       | Go 1.25                                                |
| Database       | PostgreSQL (`lib/pq`)                                  |
| Cache          | Redis (`go-redis/v9`, e.g. Upstash)                    |
| Auth           | JWT (`golang-jwt/jwt/v5`), bcrypt, Google OAuth 2.0    |
| Payment        | Midtrans (`midtrans-go`), Sandbox mode                 |
| Email          | Resend (`resend-go/v3`)                                |
| Others         | `rs/cors`, `godotenv`, `testify`, Air (hot reload)     |

---

## 📁 Project Structure

```
backend/
├── cmd/
│   └── api/
│       └── main.go            # Entry point, dependency wiring & routing
├── pkg/
│   ├── config/                # PostgreSQL & Redis connections
│   ├── cron_Handler/          # Health check endpoint for cron
│   ├── entity/                # Models / DTOs (User, Product, Cart, Order, etc.)
│   ├── handler/               # HTTP handlers (presentation layer)
│   ├── helper/                # JWT utilities
│   ├── middleware/            # JWT & RBAC middleware (+ tests)
│   ├── payment/               # Payment gateway abstraction (Midtrans)
│   ├── repository/            # Database access
│   ├── service/               # Business logic (+ tests)
│   └── utils/                 # JSON response helpers
├── migrations/                # PostgreSQL database schema
├── .github/workflows/         # GitHub Actions workflow for unit tests
├── .air.toml                  # Air hot reload configuration
├── vercel.json                # Vercel deployment configuration
└── go.mod
```

---

## 🚀 Getting Started

### Prerequisites

- [Go](https://go.dev/dl/) 1.25+
- PostgreSQL
- Redis
- Accounts for [Midtrans](https://midtrans.com/) (Sandbox), [Resend](https://resend.com/), and a Google Cloud OAuth Client

### 1. Clone & install dependencies

```bash
git clone <repo-url>
cd VELO/backend
go mod download
```

### 2. Configure the environment

Create a `.env` file in the `backend/` root folder:

```env
PORT=8080

# Database & Cache
DATABASE_URL=postgres://user:password@localhost:5432/velo?sslmode=disable
REDIS_URL=redis://localhost:6379

# JWT
SECRET_KEY=your_jwt_secret

# Midtrans
SERVER_KEY=your_midtrans_server_key

# Resend (email)
RESEND_KEY=your_resend_api_key

# Google OAuth
GOOGLE_CLIENT_ID=your_google_client_id
GOOGLE_CLIENT_SECRET=your_google_client_secret
GOOGLE_REDIRECT_URL=http://localhost:8080/api/auth/google/callback

# Frontend
FRONTEND_URL=http://localhost:3000
```

| Variable               | Description                                                  |
| ---------------------- | ------------------------------------------------------------ |
| `PORT`                 | Server port (default `8080`)                                 |
| `DATABASE_URL`         | PostgreSQL connection string                                 |
| `REDIS_URL`            | Redis URL (product cache, idempotency keys, OAuth state)     |
| `SECRET_KEY`           | Secret used to sign JWTs                                     |
| `SERVER_KEY`           | Midtrans server key                                          |
| `RESEND_KEY`           | Resend API key for verification emails                       |
| `GOOGLE_CLIENT_ID`     | Google OAuth client ID                                       |
| `GOOGLE_CLIENT_SECRET` | Google OAuth client secret                                   |
| `GOOGLE_REDIRECT_URL`  | Google OAuth callback URL                                    |
| `FRONTEND_URL`         | Frontend URL (redirect after verification / OAuth login)     |

> ⚠️ Never commit your `.env` file to the repository.

### 3. Set up the database

Create the database and apply the initial schema:

```bash
createdb velo
psql "postgres://user:password@localhost:5432/velo" -f migrations/001_init.sql
```

To create an admin account, change the role of an existing user:

```sql
UPDATE users SET role = 'admin' WHERE email = 'admin@example.com';
```

### 4. Run the server

```bash
go run ./cmd/api
```

Or with hot reload using [Air](https://github.com/air-verse/air):

```bash
air
```

The server runs at `http://localhost:8080`.

---

## 📡 API Endpoints

Legend: 🔒 = requires JWT (`jwt_token` cookie, set automatically on login), 👑 = `admin` only.

### Users

| Method | Endpoint               | Auth | Description                        |
| ------ | ---------------------- | ---- | ---------------------------------- |
| POST   | `/api/users/register`  |      | Register a new user                |
| GET    | `/api/users/verify`    |      | Verify email (`?token=`)           |
| POST   | `/api/users/login`     |      | Log in, sets the JWT cookie        |
| GET    | `/api/users/me`        | 🔒   | Get the currently logged-in user   |
| POST   | `/api/users/logout`    | 🔒   | Log out, clears the cookie         |

### OAuth

| Method | Endpoint                     | Description                         |
| ------ | ---------------------------- | ----------------------------------- |
| GET    | `/api/auth/google/login`     | Redirect to the Google login page   |
| GET    | `/api/auth/google/callback`  | Callback from Google                |

### Products

| Method | Endpoint               | Auth   | Description                        |
| ------ | ---------------------- | ------ | ---------------------------------- |
| GET    | `/api/products`        |        | List products (`?page=&limit=`)    |
| GET    | `/api/product/{id}`    |        | Product details                    |
| POST   | `/api/products`        | 🔒 👑  | Create a product                   |
| PUT    | `/api/products/{id}`   | 🔒 👑  | Update a product                   |
| DELETE | `/api/products/{id}`   | 🔒 👑  | Delete a product                   |

### Cart

| Method | Endpoint          | Auth | Description                  |
| ------ | ----------------- | ---- | ---------------------------- |
| POST   | `/api/cart`       | 🔒   | Add an item to the cart      |
| GET    | `/api/cart`       | 🔒   | View cart contents           |
| PATCH  | `/api/cart/{id}`  | 🔒   | Update item quantity         |
| DELETE | `/api/cart/{id}`  | 🔒   | Remove a single item         |
| DELETE | `/api/cart`       | 🔒   | Clear the cart               |

### Orders & Payment

| Method | Endpoint                  | Auth | Description                                                          |
| ------ | ------------------------- | ---- | -------------------------------------------------------------------- |
| POST   | `/api/checkout`           | 🔒   | Check out the cart → Midtrans (requires `X-Idempotency-Key` header)  |
| GET    | `/api/orders`             | 🔒   | User's order history                                                 |
| POST   | `/api/webhook/midtrans`   |      | Midtrans payment notification webhook                                |

### Others

| Method | Endpoint     | Description                              |
| ------ | ------------ | ---------------------------------------- |
| GET    | `/api/cron`  | Health check (for cron / keep-alive)     |

### Example Requests

```bash
# Register
curl -X POST http://localhost:8080/api/users/register \
  -H "Content-Type: application/json" \
  -d '{"name":"John","email":"john@example.com","password":"secret123"}'

# Login (saves the jwt_token cookie to cookies.txt)
curl -X POST http://localhost:8080/api/users/login \
  -c cookies.txt \
  -H "Content-Type: application/json" \
  -d '{"email":"john@example.com","password":"secret123"}'

# Add to cart
curl -X POST http://localhost:8080/api/cart \
  -b cookies.txt \
  -H "Content-Type: application/json" \
  -d '{"product_id":1,"quantity":2}'

# Checkout
curl -X POST http://localhost:8080/api/checkout \
  -b cookies.txt \
  -H "X-Idempotency-Key: 3f1c2a9e-..."
```

---

## 🔄 Order Status

| Status    | Description                                    |
| --------- | ---------------------------------------------- |
| `Unpaid`  | Order created, awaiting payment                |
| `Pending` | Payment is being processed                     |
| `Paid`    | Payment succeeded                              |
| `cancel`  | Cancelled / expired                            |

`Unpaid` orders that exceed the time limit are cancelled automatically by a scheduler that runs every 5 minutes.

---

## 🧪 Testing

```bash
go test -v ./...
```

Unit tests run automatically via **GitHub Actions** on every push / pull request to the `main` branch.

---

## 🌐 CORS

Allowed origins (see [`cmd/api/main.go`](cmd/api/main.go)):

- `http://localhost:3000`
- `https://velo-delta-gules.vercel.app`

---

## ☁️ Deployment

The project includes a `vercel.json` for deploying to **Vercel** using the `@vercel/go` runtime, with all routes directed to `api/index.go`.

> [!WARNING]
> `api/index.go` does not exist yet, so Vercel deployment will not work as-is. The current entry point is `cmd/api/main.go`, which runs a long-lived HTTP server plus a background goroutine that cancels expired orders. That goroutine will not run on a serverless platform. For serverless, the expired-order cleanup needs to be triggered by a scheduled job (e.g. Vercel Cron) instead.
