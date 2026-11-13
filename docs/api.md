# 🛒 Toko Backend API Documentation

Dokumentasi endpoint REST API resmi untuk studi kasus **Toko Backend Mentoring** (Go Fiber v3, Clean Architecture & Domain-Driven Design).

- **Base URL:** `http://localhost:8080`
- **Content-Type:** `application/json`

---

## 📌 Ringkasan Endpoint

| Method | Endpoint | Deskripsi | Status / Hands-on Task |
| :--- | :--- | :--- | :--- |
| `GET` | `/healthz` | Health check server | Standar |
| `GET` | `/api/v1/products` | Ambil seluruh daftar produk | Contoh Lengkap |
| `GET` | `/api/v1/products/:id` | Ambil detail satu produk by ID | Contoh Lengkap |
| `POST` | `/api/v1/products` | Tambah produk baru | **TODO Step 3.1 (Hands-on)** |
| `PUT` | `/api/v1/products/:id` | Update data produk by ID | **TODO Step 3.2 (Hands-on)** |
| `DELETE` | `/api/v1/products/:id` | Hapus produk by ID | **TODO Step 3.3 (Hands-on)** |
| `GET` | `/api/v1/orders` | Ambil riwayat seluruh pesanan | Contoh Lengkap (Reference) |
| `GET` | `/api/v1/orders/:id` | Ambil detail pesanan by ID | Contoh Lengkap (Reference) |
| `POST` | `/api/v1/orders` | Buat pesanan baru (Order) | Contoh Lengkap (Reference) |
| `PUT` | `/api/v1/orders/:id` | Update status/detail pesanan | Contoh Lengkap (Reference) |
| `DELETE` | `/api/v1/orders/:id` | Hapus pesanan by ID | Contoh Lengkap (Reference) |

---

## 1. System Health

### `GET /healthz`
Memeriksa apakah server Fiber v3 berjalan normal.

**cURL:**
```bash
curl -X GET http://localhost:8080/healthz
```

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "status": "healthy",
    "time": "2026-11-13T10:00:00Z",
    "version": "1.0.0"
  }
}
```

---

## 2. Products Endpoint (Fokus Hands-on CRUD)

### `GET /api/v1/products`
Mengambil semua produk yang tersedia di katalog toko.

**cURL:**
```bash
curl -X GET http://localhost:8080/api/v1/products
```

**Response (200 OK):**
```json
{
  "success": true,
  "data": [
    {
      "id": "prod-001",
      "name": "Laptop Gaming ASUS ROG",
      "price": 25000000,
      "stock": 10,
      "created_at": "2026-11-13T06:00:00Z",
      "updated_at": "2026-11-13T06:00:00Z"
    },
    {
      "id": "prod-002",
      "name": "Mechanical Keyboard",
      "price": 1200000,
      "stock": 50,
      "created_at": "2026-11-13T06:00:00Z",
      "updated_at": "2026-11-13T06:00:00Z"
    }
  ]
}
```

---

### `GET /api/v1/products/:id`
Mengambil satu produk berdasarkan ID uniknya.

**cURL:**
```bash
curl -X GET http://localhost:8080/api/v1/products/prod-001
```

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "id": "prod-001",
    "name": "Laptop Gaming ASUS ROG",
    "price": 25000000,
    "stock": 10,
    "created_at": "2026-11-13T06:00:00Z",
    "updated_at": "2026-11-13T06:00:00Z"
  }
}
```

---

### `POST /api/v1/products` *(Hands-on Step 3.1)*
Menambahkan produk baru ke katalog toko.

**Request Body Schema:**
| Field | Tipe | Wajib | Keterangan |
| :--- | :--- | :--- | :--- |
| `name` | string | Ya | Nama produk (minimal 3 karakter) |
| `price` | int64 | Ya | Harga produk (harus > 0) |
| `stock` | int | Ya | Jumlah stok awal (harus >= 0) |

**cURL:**
```bash
curl -X POST http://localhost:8080/api/v1/products \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Monitor Samsung 24 Inch",
    "price": 2300000,
    "stock": 25
  }'
```

**Response (201 Created):**
```json
{
  "success": true,
  "data": {
    "id": "c1f7b022-7728-403d-bc65-75e119d67ef9",
    "name": "Monitor Samsung 24 Inch",
    "price": 2300000,
    "stock": 25,
    "created_at": "2026-11-13T06:00:00Z",
    "updated_at": "2026-11-13T06:00:00Z"
  }
}
```

---

### `PUT /api/v1/products/:id` *(Hands-on Step 3.2)*
Memperbarui data produk berdasarkan ID.

**Request Body Schema:**
| Field | Tipe | Wajib | Keterangan |
| :--- | :--- | :--- | :--- |
| `name` | string | Ya | Nama produk baru (minimal 3 karakter) |
| `price` | int64 | Ya | Harga baru (harus > 0) |
| `stock` | int | Ya | Stok baru (harus >= 0) |

**cURL:**
```bash
curl -X PUT http://localhost:8080/api/v1/products/prod-001 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Laptop Gaming ASUS ROG Scar 16",
    "price": 27500000,
    "stock": 8
  }'
```

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "id": "prod-001",
    "name": "Laptop Gaming ASUS ROG Scar 16",
    "price": 27500000,
    "stock": 8,
    "created_at": "2026-11-13T06:00:00Z",
    "updated_at": "2026-11-13T07:15:00Z"
  }
}
```

---

### `DELETE /api/v1/products/:id` *(Hands-on Step 3.3)*
Menghapus produk dari katalog toko berdasarkan ID.

**cURL:**
```bash
curl -X DELETE http://localhost:8080/api/v1/products/prod-001
```

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "message": "Produk berhasil dihapus"
  }
}
```

---

## 3. Orders Endpoint (Contoh Lengkap / Reference)

Fitur pesanan ini berfungsi sebagai **contoh implementasi arsitektur Clean Architecture & DDD lengkap** untuk dipelajari oleh peserta.

### `GET /api/v1/orders`
Mengambil riwayat semua pesanan yang telah dibuat.

**Response (200 OK):**
```json
{
  "success": true,
  "data": [
    {
      "id": "40a5a228-56eb-46cb-8d81-8ca89a3118ee",
      "customer_id": "cust-101",
      "customer_name": "Budi Santoso",
      "items": [
        {
          "product_id": "prod-002",
          "product_name": "Mechanical Keyboard",
          "quantity": 2,
          "unit_price": 1200000,
          "subtotal": 2400000
        }
      ],
      "total": 2400000,
      "status": "PENDING",
      "created_at": "2026-11-13T06:00:00Z",
      "updated_at": "2026-11-13T06:00:00Z"
    }
  ]
}
```

### `POST /api/v1/orders`
Membuat pesanan baru dan memvalidasi aturan belanja (maksimal Rp 50.000.000).

---

## 4. Standard Error Format

Semua respon kesalahan dikelola secara terpusat oleh `CustomErrorHandler` di [`internal/middlewares/error_handler.go`](file:///D:/Kuliah/be-mentoring/internal/middlewares/error_handler.go) dengan format:

```json
{
  "success": false,
  "error": {
    "code": "<KODE_ERROR>",
    "message": "<PESAN_ERROR_DESKRIPTIF>",
    "request_id": "<UUID_REQUEST>"
  }
}
```

### Contoh Kode Kesalahan:
* **Produk Tidak Ditemukan (404 Not Found):** `PRODUCT_NOT_FOUND`
* **Nama Produk Tidak Valid (422 Unprocessable Entity):** `INVALID_PRODUCT_NAME`
* **Harga Produk Tidak Valid (422 Unprocessable Entity):** `INVALID_PRODUCT_PRICE`
* **Stok Produk Tidak Valid (422 Unprocessable Entity):** `INVALID_PRODUCT_STOCK`
* **Format JSON Rusak (400 Bad Request):** `Bad Request`
* **Limit Belanja Order Terlampaui (422 Unprocessable Entity):** `LIMIT_EXCEEDED`
