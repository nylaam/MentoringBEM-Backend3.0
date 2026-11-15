# 🛒 Toko Backend API Documentation — Product Domain

Dokumentasi REST API resmi untuk studi kasus **Toko Backend Product Service** (Go Fiber v3, Clean Architecture & Domain-Driven Design).
API ini dirancang siap deploy (Docker / Cloud PaaS) dan siap dikonsumsi langsung oleh tim **Frontend (FE)** untuk workshop atau mentoring.

- **Base URL:** `http://localhost:6767`
- **Interactive Swagger UI:** [`http://localhost:6767/docs`](http://localhost:6767/docs)
  - **Proteksi Password (HTTP Basic Auth):**
    - **Username:** `admin` *(dapat diubah di `.env`: `SWAGGER_USER`)*
    - **Password:** `admin123` *(dapat diubah di `.env`: `SWAGGER_PASSWORD`)*
- **Prefix Endpoint:** `/api/v1`
- **Content-Type:** `application/json`
- **CORS:** Diaktifkan untuk semua origin (`*`) dengan support headers `Origin, Content-Type, Accept, Authorization, X-Request-ID, X-API-Key`.
- **Autentikasi API Key untuk Endpoint Bisnis (`/api/v1/*`):**
  - Kirim melalui header: `X-API-Key: mentoring-secret-api-key-2026`
  - Atau via header: `Authorization: Bearer mentoring-secret-api-key-2026`
  - Atau via query parameter: `?api_key=mentoring-secret-api-key-2026`
- **File Spesifikasi:**
  - [`docs/openapi.json`](file:///D:/Kuliah/be-mentoring/docs/openapi.json) (OpenAPI 3.0.3 Specification)
  - [`docs/apidog_collection.json`](file:///D:/Kuliah/be-mentoring/docs/apidog_collection.json) (Collection siap pakai untuk Postman / Apidog)

---

## 📖 Cara Membuka Swagger UI di Browser

1. Buka browser dan arahkan ke alamat tunggal: **`http://localhost:6767/docs`**.
2. Browser akan memunculkan dialog popup login (**HTTP Basic Auth**).
3. Masukkan akun:
   - **Username:** `admin`
   - **Password:** `admin123`
4. Halaman interaktif Swagger UI akan tampil lengkap dengan seluruh endpoint Product CRUD & PATCH.
5. Klik tombol hijau **Authorize** di kanan atas Swagger UI, lalu masukkan `mentoring-secret-api-key-2026` pada kolom `ApiKeyAuth` untuk mencoba seluruh endpoint langsung dari browser (**Try it out**).

---

## 📌 Ringkasan Endpoint

| Method | Endpoint | Deskripsi | Auth |
| :--- | :--- | :--- | :--- |
| `GET` | `/` | Root info & status API | Public |
| `GET` | `/healthz` | Health check server & database | Public |
| `GET` | `/docs` | Interactive Swagger UI Dashboard | **Basic Auth (admin/admin123)** |
| `GET` | `/docs/openapi.json` | OpenAPI 3.0 Specification JSON | **Basic Auth (admin/admin123)** |
| `GET` | `/api/v1/products` | Ambil seluruh daftar produk | **API Key** |
| `GET` | `/api/v1/products/:id` | Ambil detail satu produk by ID | **API Key** |
| `POST` | `/api/v1/products` | Tambah produk baru (wajib: name, price, stock) | **API Key** |
| `PUT` | `/api/v1/products/:id` | Update data produk penuh (wajib: name, price, stock) | **API Key** |
| `PATCH` | `/api/v1/products/:id` | Update parsial (hanya stok, harga, atau nama) | **API Key** |
| `DELETE` | `/api/v1/products/:id` | Hapus produk by ID | **API Key** |

---

## 1. System Health & Info (Public Endpoints)

### `GET /`
Menampilkan pesan selamat datang, status, link docs, dan endpoint utama.

**cURL:**
```bash
curl -X GET http://localhost:6767/
```

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "api": "/api/v1/products",
    "docs": "/docs",
    "health": "/healthz",
    "message": "Selamat datang di Toko Backend Product API (Mentoring & Workshop)",
    "version": "1.0.0"
  }
}
```

---

### `GET /healthz`
Memeriksa apakah server Fiber v3 dan database berjalan normal (mendukung mode fallback in-memory jika PostgreSQL belum dinyalakan).

**cURL:**
```bash
curl -X GET http://localhost:6767/healthz
```

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "database": "connected",
    "status": "healthy",
    "time": "2026-11-15T15:20:00Z",
    "version": "1.0.0"
  }
}
```

---

## 2. Product Domain CRUD & PATCH (Protected with API Key)

### `GET /api/v1/products`
Mengambil semua produk yang tersedia di katalog toko.

**cURL:**
```bash
curl -X GET http://localhost:6767/api/v1/products \
  -H "X-API-Key: mentoring-secret-api-key-2026"
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
      "created_at": "2026-11-15T06:00:00Z",
      "updated_at": "2026-11-15T06:00:00Z"
    },
    {
      "id": "prod-002",
      "name": "Mechanical Keyboard",
      "price": 1200000,
      "stock": 50,
      "created_at": "2026-11-15T06:00:00Z",
      "updated_at": "2026-11-15T06:00:00Z"
    }
  ]
}
```

---

### `GET /api/v1/products/:id`
Mengambil satu produk berdasarkan ID uniknya.

**cURL:**
```bash
curl -X GET http://localhost:6767/api/v1/products/prod-001 \
  -H "X-API-Key: mentoring-secret-api-key-2026"
```

---

### `POST /api/v1/products`
Menambahkan produk baru ke katalog toko.

**Request Body Schema:**
| Field | Tipe | Wajib | Validasi | Keterangan |
| :--- | :--- | :--- | :--- | :--- |
| `name` | string | Ya | min: 3, max: 100 | Nama produk |
| `price` | int64 | Ya | gt: 0 | Harga produk dalam Rupiah |
| `stock` | int | Ya | gte: 0 | Jumlah stok awal |

**cURL:**
```bash
curl -X POST http://localhost:6767/api/v1/products \
  -H "Content-Type: application/json" \
  -H "X-API-Key: mentoring-secret-api-key-2026" \
  -d '{
    "name": "Wireless Mouse Ergonomis",
    "price": 650000,
    "stock": 40
  }'
```

**Response (201 Created):**
```json
{
  "success": true,
  "data": {
    "id": "e3e8fcbe-21bf-4cf3-a63e-6ba16db614bf",
    "name": "Wireless Mouse Ergonomis",
    "price": 650000,
    "stock": 40,
    "created_at": "2026-11-15T06:30:00Z",
    "updated_at": "2026-11-15T06:30:00Z"
  }
}
```

---

### `PUT /api/v1/products/:id` (Full Update)
Memperbarui seluruh data produk berdasarkan ID. Semua field (`name`, `price`, `stock`) wajib dikirim.

**cURL:**
```bash
curl -X PUT http://localhost:6767/api/v1/products/prod-001 \
  -H "Content-Type: application/json" \
  -H "X-API-Key: mentoring-secret-api-key-2026" \
  -d '{
    "name": "Laptop Gaming ASUS ROG Scar 16",
    "price": 27500000,
    "stock": 8
  }'
```

---

### `PATCH /api/v1/products/:id` (Partial Update)
Memperbarui data kecil/parsial pada produk. Hanya field yang ingin diubah yang perlu dikirimkan di body JSON. Minimal satu field harus disediakan.

**Contoh 1: Hanya update stok (Stock Adjustment)**
```bash
curl -X PATCH http://localhost:6767/api/v1/products/prod-001 \
  -H "Content-Type: application/json" \
  -H "X-API-Key: mentoring-secret-api-key-2026" \
  -d '{
    "stock": 25
  }'
```

**Contoh 2: Hanya ubah harga (Flash Sale / Promo)**
```bash
curl -X PATCH http://localhost:6767/api/v1/products/prod-001 \
  -H "Content-Type: application/json" \
  -H "X-API-Key: mentoring-secret-api-key-2026" \
  -d '{
    "price": 22000000
  }'
```

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "id": "prod-001",
    "name": "Laptop Gaming ASUS ROG Scar 16",
    "price": 22000000,
    "stock": 25,
    "created_at": "2026-11-15T06:00:00Z",
    "updated_at": "2026-11-15T07:30:00Z"
  }
}
```

---

### `DELETE /api/v1/products/:id`
Menghapus produk dari katalog toko berdasarkan ID.

**cURL:**
```bash
curl -X DELETE http://localhost:6767/api/v1/products/prod-001 \
  -H "X-API-Key: mentoring-secret-api-key-2026"
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

## 3. Format Respon Kesalahan Standar

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

---

## 4. Panduan Konsumsi API untuk Tim Frontend (FE)

```typescript
const BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:6767/api/v1';
const API_KEY = process.env.NEXT_PUBLIC_API_KEY || 'mentoring-secret-api-key-2026';

const defaultHeaders = {
  'Content-Type': 'application/json',
  'X-API-Key': API_KEY,
};

// 1. Fetch seluruh produk
export async function getProducts(): Promise<Product[]> {
  const res = await fetch(`${BASE_URL}/products`, {
    headers: defaultHeaders,
  });
  const json: APIResponse<Product[]> = await res.json();
  if (!json.success || !json.data) throw new Error(json.error?.message || 'Gagal memuat produk');
  return json.data;
}

// 2. Patch sebagian field (misal ubah stok)
export async function patchProductStock(id: string, newStock: number): Promise<Product> {
  const res = await fetch(`${BASE_URL}/products/${id}`, {
    method: 'PATCH',
    headers: defaultHeaders,
    body: JSON.stringify({ stock: newStock }),
  });
  const json: APIResponse<Product> = await res.json();
  if (!json.success || !json.data) throw new Error(json.error?.message || 'Gagal update stok');
  return json.data;
}
```
