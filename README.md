# Online Shop Project

## Cara Menjalankan
1. Jalankan database PostgreSQL

```
docker run --name postgresql -e POSTGRES_USER=user -e POSTGRES_PASSWORD=password -e POSTGRES_DB=database -d -p 5432:5432 postgres:16
```

2. Export Environment Variable yang dibutuhkan

```
export DB_URI=postgres://user:password@localhost:5432/database?sslmode=disable
export ADMIN_SECRET=secret
```

3. Jalankan aplikasi

```
go run .
```

## Konten

## Module
- gin: Framework web (routing, middleware, request bind & validation, etc)
- pgx: PostgreSQL driver
- uuid: ID untuk data di tabel
- crypto: Melakukan hashing passcode order/pesanan

## Route
### Publik
- [GET] /api/v1/products
- [GET] /api/v1/products/{id}
- [POST] /api/v1/checkout

### Passcode
- [POST] /api/v1/orders/{id}/confirm
- [GET] /api/v1/orders/{id}

### Admin
- [POST] /admin/products
- [PUT] /admin/products/{id}
- [DELETE] /admin/products/{id}

## Dokumentasi API
Contoh request yang memuat URL, Method, Header, dan Body dapat dilihat di folder [.http](.http)
