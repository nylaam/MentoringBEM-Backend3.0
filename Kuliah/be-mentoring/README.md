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
│   │   ├── order.go             # [HANDS-ON STEP 1.2 LOGIC]
│   │   └── product.go
│   ├── enums/                   # Status Enums (PENDING, PAID, dll.)
│   │   └── order_status.go
│   └── errx/                    # Domain Errors terisolasi per entity (Slide 13)
│       ├── errx.go              # Base DomainError & ErrorKind
│       ├── order_errx.go        # [HANDS-ON STEP 1.1 LOGIC] Error Order
│       └── product_errx.go      # Error Product
├── helpers/                     # Shared Utility Helpers
│   └── pagination/
│       └── pagination.go        # Helper Pagination (Page, Limit, Offset)
├── internal/
│   ├── app/                     # Modul Fitur Aplikasi
│   │   ├── order/
│   │   │   ├── handler.go       # [HANDS-ON STEP 3 LOGIC] Delivery Layer (Fiber v3 HTTP Handler)
│   │   │   ├── routes.go        # HTTP Route Binding
│   │   │   └── service.go       # [HANDS-ON STEP 2 LOGIC] Use Case Service
│   │   └── product/             # Fitur Produk (Contoh Implementasi Lengkap)
│   │       ├── handler.go
│   │       ├── routes.go
│   │       └── service.go
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
│       ├── order_entity_test.go # Test Table-Driven Domain Invariant
│       ├── order_service_test.go# Test Mock Repository & Service
│       ├── order_handler_test.go# Test HTTP Delivery Layer / Controller
│       └── pagination_test.go   # Test Pagination Helper
├── .env.example
├── Dockerfile                   # Multi-stage minimal production image
├── docker-compose.yml           # PostgreSQL + API container stack
├── go.mod
├── answer.md                    # 🔑 KUNCI JAWABAN LENGKAP SEMUA // TODO
└── README.md
```

---

## 🎯 Fokus Hands-on Lab: Business Logic

Pada hands-on ini, **seluruh unit test di folder `tests/unit/` sudah disediakan lengkap dan siap jalan**. Tugas peserta adalah **berfokus mengimplementasikan Logika Bisnis** di file yang bertanda **`// TODO`**:

| Step | File | Tugas Logika Bisnis |
| :--- | :--- | :--- |
| **Step 1.1** | [`domain/errx/order_errx.go`](file:///D:/Kuliah/be-mentoring/domain/errx/order_errx.go) | **Domain Error**: Definisikan `ErrOrderLimitExceeded` dengan code `"LIMIT_EXCEEDED"`. |
| **Step 1.2** | [`domain/entity/order.go`](file:///D:/Kuliah/be-mentoring/domain/entity/order.go) | **Domain Invariant**: Validasi total order tidak boleh melebihi Rp 50.000.000 (`MaxOrderTotal`). Jika melebihi, return `errx.ErrOrderLimitExceeded`. |
| **Step 2** | [`internal/app/order/service.go`](file:///D:/Kuliah/be-mentoring/internal/app/order/service.go) | **Service Orchestration**: Panggil `entity.NewOrder` dan simpan order via `s.orderRepo.Save`. |
| **Step 3** | [`internal/app/order/handler.go`](file:///D:/Kuliah/be-mentoring/internal/app/order/handler.go) | **Delivery Layer (Controller)**: Binding request body JSON, validasi request DTO, panggil service `CreateOrder`, dan return HTTP response 201 Created. |

> 💡 **Kunci Jawaban:**
> Solusi kode untuk setiap `// TODO` tersedia lengkap di file [`answer.md`](file:///D:/Kuliah/be-mentoring/answer.md).

---

## 🧪 Menguji Logika dengan Unit Test

Jalankan test suite di terminal:

```powershell
go test ./tests/unit/... -v -cover
```
*(Catatan: Tambahkan flag `-race` jika dijalankan di environment Linux atau Docker dengan CGO aktif)*

- **Sebelum `// TODO` dikerjakan:** Test invariant limit belanja dan create order akan gagal (FAIL) dan memberi petunjuk jelas apa yang belum selesai.
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

# 2. Jalankan migrasi tabel per entity secara berurutan
go run ./cmd/migrate

# (Opsional: rollback semua migrasi)
# go run ./cmd/migrate down

# 3. Masukkan data dummy produk
go run ./cmd/seed

# 4. Jalankan aplikasi
go run ./cmd/app
```
