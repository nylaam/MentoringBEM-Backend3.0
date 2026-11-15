# 🛒 Toko Backend API — Product Domain (Mentoring & Hands-on Lab)

Repositori ini adalah bahan praktik **Hands-on Mentoring & Workshop** untuk materi:
> **"Building Production-Ready Backend with Go Fiber, Clean Architecture & DDD"**

Proyek ini dibangun menggunakan **Go**, framework modern **Go Fiber v3**, arsitektur bersih **Clean Architecture & Domain-Driven Design (DDD)**, serta fokus terisolasi pada **Domain Product** (CRUD lengkap + **PATCH** untuk pembaruan cepat/kecil). API ini sudah dilengkapi **Interactive Swagger UI (Password-Protected)** di alamat tunggal `/docs`, **CORS**, validasi ketat, autentikasi **API Key**, port default **6767**, automated **CI/CD ke GitHub Container Registry (`ghcr.io`)**, dan siap dideploy langsung via Docker untuk dikonsumsi oleh tim **Frontend (FE)**.

---

## 📖 Akses Swagger UI di Browser

Swagger UI interaktif tersedia di alamat tunggal:
- **URL:** [`http://localhost:6767/docs`](http://localhost:6767/docs)
- **Autentikasi (HTTP Basic Auth):**
  - **Username:** `admin` *(konfigurasi di `.env`: `SWAGGER_USER`)*
  - **Password:** `admin123` *(konfigurasi di `.env`: `SWAGGER_PASSWORD`)*

Setelah login, Swagger UI **otomatis ter-authorize** dengan API Key dari `.env` (`mentoring-secret-api-key-2026`) sehingga Anda bisa langsung mencoba seluruh endpoint (**Try it out**).

---

## 📂 Struktur Direktori Proyek

```text
be-mentoring/
├── .github/
│   └── workflows/
│       └── ci.yml               # Automated CI/CD (Test, Build, Push GHCR & Deploy)
├── cmd/
│   ├── app/
│   │   └── main.go              # Composition Root: Dependency Injection, CORS, Swagger, API Key (Fiber v3)
│   ├── migrate/
│   │   └── main.go              # Database Migration Runner
│   └── seed/
│       └── main.go              # Database Seeder (Produk awal)
├── docs/
│   ├── api.md                   # Dokumentasi REST API lengkap, schema & cURL
│   ├── openapi.json             # OpenAPI 3.0.3 Specification (Port 6767)
│   ├── spec.go                  # Embedded OpenAPI JSON ke dalam binary (//go:embed)
│   ├── swagger.go               # Swagger UI Handler & Basic Auth Router (/docs)
│   └── apidog_collection.json   # Koleksi Postman / Apidog siap import
├── domain/                      # Core Domain Layer (Pure Go, 0 external framework)
│   ├── contracts/               # Interface Repository & Service
│   │   ├── product_repository.go
│   │   └── product_service.go
│   ├── dto/                     # Data Transfer Objects (Request & Response)
│   │   ├── product_dto.go       # Create, Update (PUT), dan Patch (PATCH) DTO
│   │   └── response_dto.go      # Standarisasi JSON Response & Error
│   ├── entity/                  # Business Logic & Domain Invariants
│   │   └── product.go           # Entity Product (NewProduct, Update, Patch)
│   └── errors/                  # Domain Errors terisolasi
│       ├── errors.go
│       └── product_errors.go
├── helpers/                     # Shared Utility Helpers
│   └── pagination/
│       └── pagination.go
├── internal/
│   ├── app/                     # Fitur Domain Aplikasi
│   │   └── product/             # Fitur Produk (CRUD & PATCH)
│   │       ├── handler.go       # Fiber v3 Thin Controller
│   │       ├── routes.go        # HTTP Route Binding (/api/v1/products)
│   │       └── service.go       # Use Case Service
│   ├── config/
│   │   └── config.go            # 12-Factor App Config (.env / API_KEY / SWAGGER_USER / Port 6767)
│   ├── infra/                   # Infrastructure Layer
│   │   ├── database/
│   │   │   └── postgres.go      # PostgreSQL Connection Pool & Auto-Table Init
│   │   ├── repository/
│   │   │   └── product_repository.go # Postgres & In-Memory Fallback Repo
│   │   └── storage/
│   │       └── storage.go
│   └── middlewares/             # Fiber v3 Middlewares (API Key, CORS, RequestID, Logger, Error Handler)
│       ├── api_key.go           # Proteksi API Key
│       ├── error_handler.go
│       ├── logger.go
│       └── request_id.go
├── migrations/                  # File DDL SQL Database Migrations
│   ├── 000001_create_products_table.up.sql
│   └── 000001_create_products_table.down.sql
├── pkg/
│   ├── response/
│   │   └── response.go          # Standardized API response builder
│   └── validator/
│       └── validator.go         # Request body payload validator (Go Playground)
├── tests/                       # Unit Tests Terpusat
│   └── unit/
│       ├── api_key_test.go        # Test Middleware API Key Auth
│       ├── pagination_test.go     # Test Pagination Helper
│       ├── product_entity_test.go # Test Invariant Domain Product (New, Update, Patch)
│       ├── product_service_test.go# Test Mock Repository & Service Product
│       ├── product_handler_test.go# Test Controller / Delivery Fiber v3
│       └── swagger_test.go        # Test Swagger UI Basic Auth & OpenAPI Endpoint
├── .env.example
├── Dockerfile                   # Multi-stage minimal production image (Port 6767)
├── docker-compose.yml           # Local Development Stack (build: .)
├── docker-compose.prod.yml      # Production VPS Stack (image: ghcr.io)
├── go.mod
├── answer.md                    # 🔑 Kunci Jawaban Lengkap Semua Tahapan Hands-on
└── README.md
```

---

## 🎯 Fitur Endpoint Product Domain (CRUD + PATCH)

- **Default Port:** `6767` (`http://localhost:6767`)
- **Autentikasi Header API:** `X-API-Key: mentoring-secret-api-key-2026`
- **Swagger Docs:** `http://localhost:6767/docs` *(Basic Auth: admin / admin123)*

| Method | Endpoint | Deskripsi | Auth |
| :--- | :--- | :--- | :--- |
| `GET` | `/` | Root Welcome API | Public |
| `GET` | `/healthz` | Health check server & database | Public |
| `GET` | `/docs` | Swagger UI Interactive API Docs | **Basic Auth** (`admin` / `admin123`) |
| `GET` | `/docs/openapi.json` | OpenAPI 3.0.3 Specification JSON | **Basic Auth** (`admin` / `admin123`) |
| `GET` | `/api/v1/products` | Ambil semua produk | **API Key** |
| `GET` | `/api/v1/products/:id` | Ambil detail produk by ID | **API Key** |
| `POST` | `/api/v1/products` | Tambah produk baru | **API Key** |
| `PUT` | `/api/v1/products/:id` | Update seluruh data produk | **API Key** |
| `PATCH`| `/api/v1/products/:id` | Update parsial produk (stok/harga/nama) | **API Key** |
| `DELETE`| `/api/v1/products/:id`| Hapus produk by ID | **API Key** |

---

## 🌐 Siap Dikonsumsi Frontend (FE Ready)

1. **CORS Otomatis Aktif:** Frontend dari `http://localhost:3000` (Next.js/React/Vite) atau domain manapun dapat langsung memanggil API tanpa isu CORS.
2. **Autentikasi Mudah:** FE cukup menambahkan header `'X-API-Key': process.env.NEXT_PUBLIC_API_KEY` pada setiap request API.
3. **Format Respon Konsisten:**
   ```json
   {
     "success": true,
     "data": { ... }
   }
   ```
4. **TypeScript Types & Contoh Fetch:** Tersedia lengkap di [`docs/api.md`](file:///D:/Kuliah/be-mentoring/docs/api.md#4-panduan-konsumsi-api-untuk-tim-frontend-fe).

---

## 🧪 Menguji Logika dengan Unit Test

Jalankan test suite terpusat:

```powershell
go test ./tests/unit/... -v -cover
```

Seluruh pengujian unit berjalan 100% (PASS) dalam hitungan milidetik tanpa perlu database aktif!

---

## 🚀 Menjalankan Server Aplikasi (Lokal)

### Cara 1: Mode Cepat Tanpa Database (Auto In-Memory)
Aplikasi otomatis mendeteksi jika PostgreSQL tidak tersedia dan langsung menggunakan **In-Memory Repository** dengan data awal di port **6767**:

```powershell
go run ./cmd/app
```
Akses di browser:
- Swagger UI: `http://localhost:6767/docs` *(Username: admin, Password: admin123)*
- Root Info: `http://localhost:6767/`
- Healthcheck: `http://localhost:6767/healthz`
- Daftar Produk: `http://localhost:6767/api/v1/products`

---

### Cara 2: Menjalankan dengan Database PostgreSQL & Docker Compose
```powershell
# 1. Nyalakan Container API (Port 6767) dan PostgreSQL di background
docker compose up -d

# 2. Periksa status container
docker compose ps
```

---

## 🚢 Panduan Deploy ke VPS via GHCR (GitHub Container Registry)

Setiap kali Anda melakukan `git push` ke branch `main`, GitHub Actions akan otomatis menguji kode, mem-build Docker Image, dan mem-publish-nya ke:
```text
ghcr.io/nylaam/mentoringbem-backend3.0:latest
```

### Cara Jalankan di VPS:

1. **Siapkan file `.env` di VPS:**
   Salin file [`.env.example`](file:///D:/Kuliah/be-mentoring/.env.example) menjadi `.env` di VPS Anda.

2. **Login ke GHCR (jika repository bersifat private):**
   ```bash
   echo $CR_PAT | docker login ghcr.io -u <github-username> --password-stdin
   ```
   *(Jika package dibuat public di GitHub Packages Settings, langkah login ini tidak diperlukan).*

3. **Tarik Image & Jalankan dengan Docker Compose:**
   ```bash
   # Jalankan menggunakan docker-compose.prod.yml
   docker compose -f docker-compose.prod.yml pull
   docker compose -f docker-compose.prod.yml up -d
   ```

4. **Atau jalankan langsung dengan Docker Run:**
   ```bash
   docker run -d \
     --name toko-api \
     --restart always \
     -p 6767:6767 \
     --env-file .env \
     ghcr.io/nylaam/mentoringbem-backend3.0:latest
   ```

Aplikasi dan Swagger UI akan langsung aktif di `http://<ip-vps-anda>:6767/docs`!
