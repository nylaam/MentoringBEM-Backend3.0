# 🛒 Toko Backend API — Mentoring Hands-on Lab

Repositori ini adalah bahan praktik **Hands-on Mentoring** untuk materi presentasi di [`index.md`](file:///D:/Kuliah/be-mentoring/index.md):
> **"Building Production-Ready Backend with Go Fiber, Clean Architecture & DDD"**

Proyek ini dibangun menggunakan **Go**, framework performan generasi terbaru **Go Fiber v3**, prinsip **Clean Architecture & Domain-Driven Design (DDD)**, serta dilengkapi **Unit Tests & Mocks** di folder `tests/unit/` dan **Docker Multi-Stage Build**.

---

## 📂 Struktur Direktori Proyek

```text
be-mentoring/
├── .github/
│   └── workflows/
│       └── ci.yml               # Automated CI (Test & Build di GitHub Actions)
├── cmd/
│   ├── app/
│   │   └── main.go              # Composition Root: Dependency Injection (Go Fiber v3)
│   ├── migrate/
│   │   └── main.go              # Database Migration Runner
│   └── seed/
│       └── main.go              # Database Seeder (Produk awal)
├── docs/
│   └── api.md                   # Dokumentasi REST API lengkap & cURL
├── domain/                      # Core Domain Layer (Pure Go, 0 external framework)
│   ├── contracts/               # Interface Repository & Service
│   │   ├── order_repository.go
│   │   ├── order_service.go
│   │   ├── product_repository.go
│   │   └── product_service.go
│   ├── dto/                     # Data Transfer Objects (Request & Response)
│   │   ├── order_dto.go
│   │   ├── product_dto.go
│   │   └── response_dto.go
│   ├── entity/                  # Business Logic & Invariants
│   │   ├── order.go             # Entity Order (Contoh Implementasi Lengkap)
│   │   └── product.go           # [HANDS-ON STEP 1 LOGIC] Invariant Produk
│   ├── enums/                   # Status Enums (PENDING, PAID, dll.)
│   │   └── order_status.go
│   └── errors/                  # Domain Errors terisolasi per entity (Slide 13)
│       ├── errors.go            # Base DomainError & ErrorKind
│       ├── order_errors.go      # Error Order
│       └── product_errors.go    # Error Produk
├── helpers/                     # Shared Utility Helpers
│   └── pagination/
│       └── pagination.go        # Helper Pagination (Page, Limit, Offset)
├── internal/
│   ├── app/                     # Modul Fitur Aplikasi
│   │   ├── order/               # Fitur Order (Contoh Implementasi Lengkap / Reference)
│   │   │   ├── handler.go       # Fiber v3 HTTP Handler
│   │   │   ├── routes.go        # HTTP Route Binding
│   │   │   └── service.go       # Use Case Service
│   │   └── product/             # Fitur Produk [HANDS-ON CRUD UTAMA]
│   │       ├── handler.go       # [HANDS-ON STEP 3 LOGIC] Delivery Layer
│   │       ├── routes.go        # HTTP Route Binding
│   │       └── service.go       # [HANDS-ON STEP 2 LOGIC] Use Case Service
│   ├── config/
│   │   └── config.go            # 12-Factor App Configuration (.env)
│   ├── infra/                   # Infrastructure Layer
│   │   ├── database/
│   │   │   └── postgres.go      # PostgreSQL Connection Pool
│   │   ├── repository/
│   │   │   ├── order_repository.go
│   │   │   └── product_repository.go
│   │   ├── scheduler/           # Background Job Scheduler
│   │   │   └── scheduler.go
│   │   └── storage/             # File Storage Abstraction
│   │       └── storage.go
│   └── middlewares/             # Fiber v3 Middlewares (RequestID, Logger, Error Handler)
│       ├── error_handler.go
│       ├── logger.go
│       └── request_id.go
├── migrations/                  # File DDL SQL Database Migrations per Entity
│   ├── 000001_create_products_table.up.sql
│   ├── 000001_create_products_table.down.sql
│   ├── 000002_create_orders_table.up.sql
│   ├── 000002_create_orders_table.down.sql
│   ├── 000003_create_order_items_table.up.sql
│   └── 000003_create_order_items_table.down.sql
├── pkg/
│   ├── response/
│   │   └── response.go          # Standardized API response builder
│   └── validator/
│       └── validator.go         # Request body payload validator
├── tests/                       # Semua File Test Terpusat di tests/unit
│   └── unit/
│       ├── product_entity_test.go # [HANDS-ON] Test Domain Invariant Produk
│       ├── product_service_test.go# [HANDS-ON] Test Mock Repository & Service Produk
│       ├── product_handler_test.go# [HANDS-ON] Test Controller / Delivery Produk
│       ├── order_entity_test.go   # Test Reference Order Entity
│       ├── order_service_test.go  # Test Reference Order Service
│       ├── order_handler_test.go  # Test Reference Order Handler
│       └── pagination_test.go     # Test Pagination Helper
├── .env.example
├── Dockerfile                   # Multi-stage minimal production image
├── docker-compose.yml           # PostgreSQL + API container stack
├── go.mod
├── answer.md                    # 🔑 KUNCI JAWABAN LENGKAP SEMUA // TODO
└── README.md
```

---

## 🎯 Fokus Hands-on Lab: Product CRUD

Hands-on mentoring ini berfokus pada **Product CRUD** dengan alur Clean Architecture (Domain Entity ➔ Service Use Case ➔ Delivery Controller). Fitur Order disediakan sebagai contoh implementasi lengkap (*reference architecture*).

Tugas peserta adalah menyelesaikan bagian bertanda **`// TODO`**:

| Step | File | Tugas Hands-on Product CRUD |
| :--- | :--- | :--- |
| **Step 1.1** | [`domain/entity/product.go`](file:///D:/Kuliah/be-mentoring/domain/entity/product.go) | **Domain Invariant `NewProduct`**: Validasi nama minimal 3 karakter, harga > 0, dan stok >= 0. |
| **Step 1.2** | [`domain/entity/product.go`](file:///D:/Kuliah/be-mentoring/domain/entity/product.go) | **Domain Invariant `Update`**: Validasi perubahan nama, harga, stok, serta perbarui `UpdatedAt`. |
| **Step 2.1** | [`internal/app/product/service.go`](file:///D:/Kuliah/be-mentoring/internal/app/product/service.go) | **Use Case `CreateProduct`**: Panggil `entity.NewProduct` dan simpan via `s.repo.Save`. |
| **Step 2.2** | [`internal/app/product/service.go`](file:///D:/Kuliah/be-mentoring/internal/app/product/service.go) | **Use Case `UpdateProduct`**: Ambil produk via `s.repo.FindByID`, jalankan `p.Update(...)`, simpan via `s.repo.Update`. |
| **Step 2.3** | [`internal/app/product/service.go`](file:///D:/Kuliah/be-mentoring/internal/app/product/service.go) | **Use Case `DeleteProduct`**: Cek keberadaan produk dan hapus via `s.repo.Delete`. |
| **Step 3.1** | [`internal/app/product/handler.go`](file:///D:/Kuliah/be-mentoring/internal/app/product/handler.go) | **Controller (POST /products)**: Binding JSON, validasi request DTO, panggil service `CreateProduct`, return 201 Created. |
| **Step 3.2** | [`internal/app/product/handler.go`](file:///D:/Kuliah/be-mentoring/internal/app/product/handler.go) | **Controller (PUT /products/:id)**: Ambil param `id`, bind & validasi `UpdateProductRequest`, panggil `UpdateProduct`, return 200 OK. |
| **Step 3.3** | [`internal/app/product/handler.go`](file:///D:/Kuliah/be-mentoring/internal/app/product/handler.go) | **Controller (DELETE /products/:id)**: Ambil param `id`, panggil `DeleteProduct`, return 200 OK dengan pesan sukses. |

---

## 🧪 Menguji Logika dengan Unit Test

Jalankan test suite di terminal:

```powershell
go test ./tests/unit/... -v -cover
```
*(Catatan: Tambahkan flag `-race` jika dijalankan di environment Linux atau Docker dengan CGO aktif)*

- **Sebelum `// TODO` dikerjakan:** Test Product CRUD akan gagal (FAIL) dan memberi petunjuk jelas apa yang belum selesai.
- **Setelah `// TODO` diselesaikan:** Seluruh pengujian akan lulus 100% (PASS) dalam hitungan milidetik tanpa membutuhkan database!

---

## 🚀 Menjalankan Server Aplikasi

### Cara 1: Mode Cepat Tanpa Database (Auto In-Memory)
Aplikasi otomatis mendeteksi jika Postgres belum aktif dan beralih ke mode **In-Memory Repository**:

```powershell
go run ./cmd/app
```
Akses di browser / Postman / cURL:
- Healthcheck: `http://localhost:8080/healthz`
- Daftar Produk: `http://localhost:8080/api/v1/products`
- Dokumentasi Lengkap: [`docs/api.md`](file:///D:/Kuliah/be-mentoring/docs/api.md)

---

### Cara 2: Menjalankan dengan Database PostgreSQL & Docker Compose
```powershell
# 1. Nyalakan PostgreSQL
docker compose up -d

# 2. Jalankan migrasi tabel
go run ./cmd/migrate

# (Opsional: rollback semua migrasi)
# go run ./cmd/migrate down

# 3. Masukkan data dummy produk
go run ./cmd/seed

# 4. Jalankan aplikasi
go run ./cmd/app
```
