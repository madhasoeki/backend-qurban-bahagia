# API Contract - Qurban Service

Base URL: /api

## Authentication
- Header: Authorization: Bearer <token>
- Token expires in 24 hours.
- Roles: admin, koordinator_pengawas, pengawas, jagal, kulit, cacah_daging, cacah_tulang, packing, distribusi

## Common Responses
- Success: JSON with keys like message, data
- Error: { "error": "<message>" }

## Public Endpoints

### POST /login
Authenticate user and return JWT.

Request body:
```json
{
  "username": "string",
  "password": "string"
}
```

Response 200:
```json
{
  "username": "string",
  "token": "string",
  "role": "admin"
}
```

Errors:
- 400 Input tidak valid
- 401 Username atau password salah
- 500 Gagal membuat token

### GET /stream
Server-Sent Events (SSE) stream.

Headers:
- Content-Type: text/event-stream

Event payload (example):
```json
{
  "action": "UPDATE_HEWAN",
  "data": { }
}
```

Actions:
- UPDATE_HEWAN
- UPDATE_DASHBOARD
- UPDATE_DISTRIBUSI

### GET /dashboard/summary
Public dashboard summary.

Response 200:
```json
{
  "data": {
    "total_hewan": 0,
    "total_hewan_selesai": 0,
    "total_kantong_packing": 0,
    "total_kantong_distribusi": 0,
    "waktu_mulai": "2026-05-08T10:00:00Z",
    "waktu_selesai": "2026-05-08T15:00:00Z"
  }
}
```

### GET /dashboard/hewan
Public list of hewan (includes pengawas).

Response 200:
```json
{
  "message": "Berhasil",
  "data": [
    {
      "id": 1,
      "kode_hewan": "Q001",
      "tipe": "qurban",
      "jenis_hewan": "sapi",
      "nama_sohibul": ["Nama A"],
      "catatan": "string",
      "waktu_mulai_jagal": null,
      "waktu_selesai_jagal": null,
      "waktu_mulai_kuliti": null,
      "waktu_selesai_kuliti": null,
      "waktu_mulai_cacah_daging": null,
      "waktu_selesai_cacah_daging": null,
      "waktu_mulai_cacah_tulang": null,
      "waktu_selesai_cacah_tulang": null,
      "waktu_mulai_packing": null,
      "waktu_selesai_packing": null,
      "kantong_packing": null,
      "berat_daging": 0,
      "berat_tulang": 0,
      "cek_kepala": false,
      "cek_kaki": false,
      "cek_kulit": false,
      "cek_ekor": false,
      "cek_distribusi": false,
      "pengawas_id": 2,
      "pengawas": {
        "id": 2,
        "nama_lengkap": "string",
        "username": "string",
        "role": "pengawas"
      },
      "created_at": "2026-05-08T10:00:00Z",
      "updated_at": "2026-05-08T10:00:00Z"
    }
  ]
}
```

## Authenticated Endpoints (All Roles)

### GET /hewan
Requires any role in AuthMiddleware.

Query params:
- search: string (matches kode_hewan, nama_sohibul, catatan)
- tipe: qurban | sedekah
- jenis_hewan: sapi | kambing
- pengawas_id: number (ignored for role pengawas, forced to own ID)

Response 200: same shape as /dashboard/hewan

## Operational Endpoints (admin, koordinator_pengawas, pengawas)

### PATCH /hewan/:id/progress/:pos
Update process timestamps. pos: jagal | kulit | cacah_daging | cacah_tulang.

Request body:
```json
{ "status": "mulai" }
```

Rules:
- For pos != jagal, jagal must be completed
- For role pengawas: can only update own hewan

Response 200:
```json
{
  "message": "Progress jagal berhasil diperbarui",
  "data": { }
}
```

Errors:
- 400 Status harus 'mulai' atau 'selesai'
- 400 Proses jagal belum selesai
- 400 Pos operasional tidak dikenali
- 403 Anda tidak berhak mengelola hewan ini
- 404 Hewan tidak ditemukan

### PATCH /hewan/:id/timbang
Update berat_daging and berat_tulang. Jagal must be completed.

Request body:
```json
{
  "berat_daging": 10.5,
  "berat_tulang": 3.2
}
```

Response 200:
```json
{
  "message": "Data timbangan berhasil disimpan",
  "data": { }
}
```

### PATCH /hewan/:id/kelengkapan
Update check flags. For role pengawas: only own hewan.

Request body (all optional):
```json
{
  "cek_kepala": true,
  "cek_kaki": true,
  "cek_kulit": false,
  "cek_ekor": false,
  "cek_distribusi": true
}
```

Response 200:
```json
{
  "message": "Data kelengkapan berhasil diperbarui",
  "data": { }
}
```

## Packing Endpoints (admin, packing)

### PATCH /hewan/:id/packing
Start or finish packing.

Request body:
```json
{ "status": "mulai" }
```

Finish packing:
```json
{ "status": "selesai", "total_kantong": 12 }
```

Rules:
- Jagal must be completed
- total_kantong required when status=selesai

Response 200:
```json
{
  "message": "Data packing berhasil diperbarui",
  "data": { }
}
```

## Distribution Endpoints (admin, distribusi)

### GET /distribusi
List distribusi with user.

Response 200:
```json
{
  "message": "Berhasil",
  "data": [
    {
      "id": 1,
      "user_id": 5,
      "jumlah_kantong": 10,
      "user": {
        "id": 5,
        "nama_lengkap": "string",
        "username": "string",
        "role": "distribusi"
      },
      "created_at": "2026-05-08T10:00:00Z",
      "updated_at": "2026-05-08T10:00:00Z"
    }
  ]
}
```

### GET /distribusi/:user_id
Get distribusi for a user. If no data, returns placeholder with jumlah_kantong=0.

Response 200 (no record):
```json
{
  "message": "Belum ada distribusi",
  "data": { "JumlahKantong": 0, "UserID": "5" }
}
```

### PATCH /distribusi/:user_id
Increment or decrement jumlah_kantong.

Request body:
```json
{ "penambahan": 3 }
```

Rules:
- Role distribusi can only update own user_id
- jumlah_kantong cannot be below 0

Response 200:
```json
{
  "message": "Data distribusi berhasil diupdate",
  "data": { }
}
```

## Admin Endpoints (admin only)

### GET /users
Query params:
- role: role string
- search: username substring

Response 200:
```json
{ "message": "Berhasil", "data": [ ] }
```

### POST /users
Request body:
```json
{
  "nama_lengkap": "string",
  "username": "string",
  "password": "string",
  "role": "pengawas"
}
```

Response 201:
```json
{ "message": "User berhasil didaftarkan", "user_id": 1 }
```

### PUT /users/:id
Request body:
```json
{
  "username": "string",
  "password": "string",
  "role": "admin"
}
```

Response 200:
```json
{ "message": "User berhasil diperbarui" }
```

### DELETE /users/:id
Response 200:
```json
{ "message": "User berhasil dihapus" }
```

### POST /hewan
Request body:
```json
{
  "kode_hewan": "Q001",
  "tipe": "qurban",
  "jenis_hewan": "sapi",
  "nama_sohibul": ["Nama A"],
  "pengawas_id": 2,
  "catatan": "string"
}
```

Response 201:
```json
{ "message": "Data hewan berhasil diinput", "data": { } }
```

### PUT /hewan/:id

Request body (same as POST /hewan, all optional):
```json
{
  "kode_hewan": "Q001",
  "tipe": "qurban",
  "jenis_hewan": "sapi",
  "nama_sohibul": ["Nama A"],
  "pengawas_id": 2,
  "catatan": "string"
}
```

Response 200:
```json
{ "message": "Data hewan berhasil diperbarui" }
```

### DELETE /hewan/:id
Rule: cannot delete if waktu_mulai_jagal is set.

Response 200:
```json
{ "message": "Data hewan berhasil dihapus" }
```
