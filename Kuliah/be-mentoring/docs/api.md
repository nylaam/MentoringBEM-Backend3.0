# Toko Backend API Documentation

Dokumentasi endpoint REST API untuk studi kasus Toko Backend (Hands-on Mentoring).

Base URL: `http://localhost:8080`

---

## 1. Health Check

### `GET /healthz`
Memeriksa status kesehatan server.

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "status": "healthy",
    "time": "2026-11-11T06:00:00Z",
    "version": "1.0.0"
  }
}
```

---

## 2. Products Endpoint

### `GET /api/v1/products`
Mengambil seluruh daftar produk di toko.

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
      "created_at": "2026-11-11T06:00:00Z",
      "updated_at": "2026-11-11T06:00:00Z"
    }
  ]
}
```

### `GET /api/v1/products/:id`
Mengambil detail produk berdasarkan ID.

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "id": "prod-001",
    "name": "Laptop Gaming ASUS ROG",
    "price": 25000000,
    "stock": 10,
    "created_at": "2026-11-11T06:00:00Z",
    "updated_at": "2026-11-11T06:00:00Z"
  }
}
```

### `POST /api/v1/products`
Menambahkan produk baru.

**Request Body:**
```json
{
  "name": "Monitor Samsung 24 Inch",
  "price": 2300000,
  "stock": 25
}
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
    "created_at": "2026-11-11T06:00:00Z",
    "updated_at": "2026-11-11T06:00:00Z"
  }
}
```

---

## 3. Orders Endpoint

### `GET /api/v1/orders`
Mengambil riwayat pesanan (orders).

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
      "created_at": "2026-11-11T06:00:00Z",
      "updated_at": "2026-11-11T06:00:00Z"
    }
  ]
}
```

### `POST /api/v1/orders`
Membuat pesanan baru (Order).

**Request Body (Contoh Sukses):**
```json
{
  "customer_id": "cust-101",
  "customer_name": "Budi Santoso",
  "items": [
    {
      "product_id": "prod-001",
      "quantity": 1
    },
    {
      "product_id": "prod-002",
      "quantity": 2
    }
  ]
}
```

**Response Sukses (201 Created):**
```json
{
  "success": true,
  "data": {
    "id": "40a5a228-56eb-46cb-8d81-8ca89a3118ee",
    "customer_id": "cust-101",
    "customer_name": "Budi Santoso",
    "items": [
      {
        "product_id": "prod-001",
        "product_name": "Laptop Gaming ASUS ROG",
        "quantity": 1,
        "unit_price": 25000000,
        "subtotal": 25000000
      },
      {
        "product_id": "prod-002",
        "product_name": "Mechanical Keyboard",
        "quantity": 2,
        "unit_price": 1200000,
        "subtotal": 2400000
      }
    ],
    "total": 27400000,
    "status": "PENDING",
    "created_at": "2026-11-11T06:00:00Z",
    "updated_at": "2026-11-11T06:00:00Z"
  }
}
```

### Error Responses

#### 1. Batas Maksimal Belanja Terlampaui (> Rp 50.000.000)
**Response (422 Unprocessable Entity):**
```json
{
  "success": false,
  "error": {
    "code": "LIMIT_EXCEEDED",
    "message": "total order melebihi batas maksimal",
    "request_id": "8fa134d4-c36b-4560-844c-00bd213605d3"
  }
}
```

#### 2. Produk Tidak Ditemukan
**Response (404 Not Found):**
```json
{
  "success": false,
  "error": {
    "code": "PRODUCT_NOT_FOUND",
    "message": "produk tidak ditemukan",
    "request_id": "8fa134d4-c36b-4560-844c-00bd213605d3"
  }
}
```

#### 3. Format Request JSON Salah / Validasi Gagal
**Response (422 Unprocessable Entity):**
```json
{
  "success": false,
  "error": {
    "code": "Unprocessable Entity",
    "message": "field 'Items' gagal validasi 'required'",
    "request_id": "8fa134d4-c36b-4560-844c-00bd213605d3"
  }
}
```
