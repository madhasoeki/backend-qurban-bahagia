# 🐄 Qurban Bahagia — Backend API

Aplikasi untuk memantau proses pemotongan hewan qurban secara real-time untuk agenda **Qurban Bahagia** di **Masjid Ismuhu Yahya**. Dibangun dengan **Go**, **Gin**, dan **GORM** sebagai backend REST API yang mendukung alur kerja mulai dari penyembelihan hingga distribusi daging.

## ✨ Fitur Utama

- **Autentikasi JWT** — Login berbasis token dengan role-based access control (RBAC)
- **Manajemen Hewan** — CRUD data hewan qurban/sedekah dengan pelacakan status per tahap
- **Alur Operasional Bertahap** — Tracking progress: Jagal → Kuliti → Cacah Daging → Cacah Tulang → Packing
- **Penimbangan** — Pencatatan berat daging dan tulang
- **Checklist Kelengkapan** — Verifikasi kepala, kaki, kulit, ekor, dan distribusi
- **Packing** — Pencatatan jumlah kantong packing per hewan
- **Distribusi** — Manajemen distribusi kantong per petugas
- **Dashboard Publik** — Ringkasan statistik dan daftar hewan tanpa autentikasi
- **Real-time Updates** — Server-Sent Events (SSE) untuk push update ke client
- **Auto Migration & Seeder** — Database otomatis di-migrate dan admin default di-seed saat startup

## 🏗️ Arsitektur

```
qurban/
├── main.go                  # Entry point, CORS config, server startup
├── config/
│   ├── database.go          # Koneksi database MySQL via GORM
│   └── seeder.go            # Seeder admin user default
├── controllers/
│   ├── auth_controller.go       # Login & JWT generation
│   ├── admin_controller.go      # CRUD user (admin only)
│   ├── hewan_controller.go      # CRUD hewan (admin) & list (all roles)
│   ├── pos_controller.go        # Progress, timbang, kelengkapan, packing
│   ├── distribusi_controller.go # Distribusi CRUD
│   ├── dashboard_controller.go  # Dashboard publik
│   └── sse_controller.go        # Server-Sent Events handler
├── middleware/
│   └── auth_middleware.go   # JWT verification & role authorization
├── models/
│   ├── user.go              # User model dengan bcrypt password hashing
│   ├── hewan.go             # Hewan model (semua field operasional)
│   ├── distribusi.go        # Distribusi model
│   └── dashboard_summary.go # Struct untuk response dashboard
├── routes/
│   └── routes.go            # Route grouping berdasarkan role
├── utils/
│   └── jwt.go               # JWT token generation & parsing
├── api.md                   # API contract documentation
├── .env.example             # Template environment variables
├── go.mod
└── go.sum
```

## 🛠️ Tech Stack

| Komponen | Teknologi |
|----------|-----------|
| Bahasa | Go 1.25 |
| Framework | Gin v1.12 |
| ORM | GORM v1.31 |
| Database | MySQL |
| Auth | JWT (`golang-jwt/jwt/v5`) |
| Password | bcrypt (`golang.org/x/crypto`) |
| CORS | `gin-contrib/cors` |
| Env | `godotenv` |

## 🚀 Getting Started

### Prerequisites

- [Go](https://go.dev/dl/) 1.25+
- [MySQL](https://dev.mysql.com/downloads/) 8.0+

### 1. Clone Repository

```bash
git clone https://github.com/<username>/qurban.git
cd qurban
```

### 2. Setup Environment Variables

Salin file `.env.example` dan isi dengan konfigurasi Anda:

```bash
cp .env.example .env
```

Isi file `.env`:

```env
APP_PORT=8080
GIN_MODE=debug

DB_HOST=127.0.0.1
DB_PORT=3306
DB_USER=root
DB_PASSWORD=your_password
DB_NAME=qurban_db

JWT_SECRET=ganti-dengan-secret-yang-kuat-minimal-32-karakter
ADMIN_DEFAULT_PASSWORD=your_admin_password
CORS_ORIGIN=http://localhost:5173
```

> **⚠️ Penting:** `JWT_SECRET` harus berupa string yang kuat (minimal 32 karakter). `ADMIN_DEFAULT_PASSWORD` wajib diisi untuk seeder admin.

### 3. Buat Database

```sql
CREATE DATABASE qurban_db;
```

### 4. Jalankan Server

```bash
go mod download
go run main.go
```

Server akan berjalan di `http://localhost:8080` (atau sesuai `APP_PORT`). Database akan otomatis di-migrate dan admin user akan di-seed pada startup pertama.

## 👥 Roles & Permissions

| Role | Akses |
|------|-------|
| `admin` | Full access — CRUD user, hewan, semua operasional |
| `koordinator_pengawas` | Progress, timbang, kelengkapan |
| `pengawas` | Progress, timbang, kelengkapan (hanya hewan sendiri) |
| `jagal` | Baca data hewan |
| `kulit` | Baca data hewan |
| `cacah_daging` | Baca data hewan |
| `cacah_tulang` | Baca data hewan |
| `packing` | Baca data hewan, update packing |
| `distribusi` | Baca data hewan, manage distribusi (hanya data sendiri) |

## 📡 API Endpoints

### Public (Tanpa Autentikasi)

| Method | Endpoint | Deskripsi |
|--------|----------|-----------|
| `POST` | `/api/login` | Login dan dapatkan JWT token |
| `GET` | `/api/stream` | SSE stream untuk real-time updates |
| `GET` | `/api/dashboard/summary` | Ringkasan dashboard |
| `GET` | `/api/dashboard/hewan` | Daftar hewan publik |

### Authenticated (Semua Role)

| Method | Endpoint | Deskripsi |
|--------|----------|-----------|
| `GET` | `/api/hewan` | Daftar hewan (filter: `search`, `tipe`, `jenis_hewan`, `pengawas_id`) |

### Operasional (admin, koordinator_pengawas, pengawas)

| Method | Endpoint | Deskripsi |
|--------|----------|-----------|
| `PATCH` | `/api/hewan/:id/progress/:pos` | Update progress per pos (jagal, kulit, cacah_daging, cacah_tulang) |
| `PATCH` | `/api/hewan/:id/timbang` | Update berat daging & tulang |
| `PATCH` | `/api/hewan/:id/kelengkapan` | Update checklist kelengkapan |

### Packing (admin, packing)

| Method | Endpoint | Deskripsi |
|--------|----------|-----------|
| `PATCH` | `/api/hewan/:id/packing` | Mulai/selesai packing (+ total kantong) |

### Distribusi (admin, distribusi)

| Method | Endpoint | Deskripsi |
|--------|----------|-----------|
| `GET` | `/api/distribusi` | Daftar semua distribusi |
| `GET` | `/api/distribusi/:user_id` | Detail distribusi per user |
| `PATCH` | `/api/distribusi/:user_id` | Update jumlah kantong (increment/decrement) |

### Admin Only

| Method | Endpoint | Deskripsi |
|--------|----------|-----------|
| `GET` | `/api/users` | Daftar user (filter: `role`, `search`) |
| `POST` | `/api/users` | Tambah user baru |
| `PUT` | `/api/users/:id` | Edit user |
| `DELETE` | `/api/users/:id` | Hapus user |
| `POST` | `/api/hewan` | Tambah hewan |
| `PUT` | `/api/hewan/:id` | Edit hewan |
| `DELETE` | `/api/hewan/:id` | Hapus hewan (gagal jika sudah diproses) |

> 📖 Dokumentasi API lengkap tersedia di file [`api.md`](api.md).

## 🔄 Alur Operasional

```
┌─────────┐    ┌─────────┐    ┌──────────────┐    ┌──────────────┐    ┌─────────┐    ┌───────────┐
│  Jagal   │───▶│  Kuliti  │───▶│ Cacah Daging │───▶│ Cacah Tulang │───▶│ Packing │───▶│Distribusi │
│ (mulai/  │    │ (mulai/  │    │   (mulai/    │    │   (mulai/    │    │ (mulai/ │    │(increment/│
│ selesai) │    │ selesai) │    │   selesai)   │    │   selesai)   │    │selesai) │    │decrement) │
└─────────┘    └─────────┘    └──────────────┘    └──────────────┘    └─────────┘    └───────────┘
                                                                           │
                                                                     ┌─────┴──────┐
                                                                     │  Timbang   │
                                                                     │(berat daging│
                                                                     │& tulang)    │
                                                                     └────────────┘
```

**Aturan:**
- Setiap pos hanya bisa dimulai jika proses **jagal** sudah selesai
- Setiap proses memiliki status `mulai` dan `selesai`
- Packing membutuhkan `total_kantong` saat status `selesai`
- Hewan yang sudah diproses (waktu_mulai_jagal terisi) tidak bisa dihapus

## 🔐 Keamanan

- Password di-hash menggunakan **bcrypt** sebelum disimpan
- Token JWT expire dalam **24 jam**
- CORS dikonfigurasi per environment via `CORS_ORIGIN`
- Role-based middleware memvalidasi akses di setiap endpoint
- Pengawas hanya bisa mengakses hewan yang ditugaskan kepadanya

## 🌐 Environment Variables

| Variable | Deskripsi | Default |
|----------|-----------|---------|
| `APP_PORT` | Port server | `8080` |
| `GIN_MODE` | Mode Gin (`debug`/`release`) | `debug` |
| `DB_HOST` | Host database MySQL | - |
| `DB_PORT` | Port database MySQL | - |
| `DB_USER` | Username database | - |
| `DB_PASSWORD` | Password database | - |
| `DB_NAME` | Nama database | - |
| `JWT_SECRET` | Secret key untuk JWT signing | - |
| `ADMIN_DEFAULT_PASSWORD` | Password admin default (seeder) | - |
| `CORS_ORIGIN` | Allowed CORS origin | `http://localhost:5173` |

## 📄 Lisensi

Proyek ini dilisensikan di bawah [MIT License](LICENSE).
