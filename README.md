# VELO Backend API

VELO Backend adalah layanan RESTful API untuk platform e-commerce (mencakup produk, keranjang, pesanan, dan pengguna) yang dibangun dengan bahasa pemrograman [Go (Golang)](https://go.dev/). Proyek ini dirancang untuk dapat di-_deploy_ dalam lingkungan _Serverless_ seperti [Render](https://render.com).

## 🚀 Fitur Utama

- **User Management**: Registrasi, Login otomatis dengan JSON Web Token (JWT).
- **Products**: Mengelola katalog produk.
- **Cart & Order**: Manajemen keranjang belanja dan siklus hidup pesanan (checkout).
- **Payment Gateway**: Integrasi pembayaran otomatis menggunakan **Midtrans**.
- **Role-Based Access Control (RBAC)**: Middleware keamanan yang membatasi akses berdasarkan _role_ pengguna.
- **Email Service**: Integrasi pengiriman email transaksional menggunakan **Resend**.

## 🛠️ Tech Stack

- **Languange:** Go (1.25+)
- **Database:** PostgreSQL (`github.com/lib/pq`)
- **Cache**: Redis (`github.com/redis/go-redis`)
- **Authentication:** JWT (`github.com/golang-jwt/jwt/v5`) & Crypto (`golang.org/x/crypto`)
- **Payment Gateway:** Midtrans (`github.com/midtrans/midtrans-go`)
- **Email Provider:** Resend (`github.com/resend/resend-go`)
- **Environment:** `godotenv`

## 📁 Struktur Folder

```text
.
├── api/             # Entrypoint alternatif (jika menggunakan serverless)
├── cmd/             # Entrypoint untuk menjalankan aplikasi (main.go)
│   └── api/
├── pkg/
│   ├── config/      # Konfigurasi Database (PostgreSQL) & Redis
│   ├── cron_Handler/# Handler dan logika untuk background jobs
│   ├── entity/      # Tipe data & struct (Cart, Order, Product, User, dll)
│   ├── handler/     # Controller untuk HTTP Request (Cart, Order, Product)
│   ├── helper/      # Utility helpers (contoh: JWT Token generator)
│   ├── middleware/  # Middleware API (JWT Auth, RBAC)
│   ├── payment/     # Modul Interaksi dengan Payment Gateway (Midtrans)
│   ├── repository/  # Query komunikasi langsung dengan PostgreSQL
│   ├── service/     # Cangkupan Business logic
│   └── utils/       # Utility Response builder
├── vercel.json      # (Opsional) Konfigurasi jika deployment menggunakan Vercel
├── go.mod           # Dependency management
└── go.sum
```

## ⚙️ Persyaratan (Prerequisites)

- **Go** v1.25 atau terbaru
- **PostgreSQL** Database aktif
- **Redis** server aktif
- Akun **Midtrans** (untuk testing sandbox / production)
- Akun **Resend** (untuk API key email)

## 📦 Instalasi & Cara Menjalankan (Local)

1. **Clone repositori ini:**

   ```bash
   git clone https://github.com/Onlyadmirer/API-VELO.git
   cd VELO-backend
   ```

2. **Download seluruh _dependencies_:**

   ```bash
   go mod tidy
   ```

3. **Duplikat/Buat file konfigurasi environment:**
   Buat file `.env` pada _root folder_, dan isi sesuai pengaturan sistem Anda:

   ```env
   # SERVER
   PORT=8080

   # DATABASE (POSTGRESQL)
   DATABASE_URL=postgres://user:pass@localhost:5432/velo_db?sslmode=disable

   # REDIS
   REDIS_URL=redis://localhost:6379/0

   # JWT
   JWT_SECRET=rahasia123

   # MIDTRANS
   MIDTRANS_SERVER_KEY=SB-Mid-server-xxxx

   # RESEND EMAIL
   RESEND_API_KEY=re_xxxx...
   ```

4. **Jalankan Aplikasi:**
   ```bash
   go run cmd/api/main.go
   ```
   atau
   ```bash
   air
   ```

Aplikasi akan berjalan secara default di port (e.g., `http://localhost:8080`).

## ☁️ Deployment (Render)

Proyek ini dapat dengan mudah di-_deploy_ sebagai **Web Service** di Render:

1. Buat **Web Service** baru di dashboard [Render](https://render.com).
2. Hubungkan repositori GitHub proyek ini.
3. Gunakan konfigurasi berikut pada pengaturan service:
   - **Environment**: `Go`
   - **Build Command**: `go build -o api-velo ./cmd/api/main.go`
   - **Start Command**: `./api-velo`
4. Tambahkan semua _Environment Variables_ yang ada di file `.env` ke bagian **Environment** di dashboard Render Anda.
5. Simpan dan tunggu proses _build_ serta _deployment_ selesai.
