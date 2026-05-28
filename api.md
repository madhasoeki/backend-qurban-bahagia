# API Contract - Qurban Service

Base URL: /api

## Authentication
- Header: Authorization: Bearer <token>
- Token expires in 24 hours.
- Roles: admin, koordinator_pengawas, pengawas, jagal, kulit, cacah_daging, cacah_tulang, packing, distribusi

## Common Responses
- Success: JSON with keys like message, data
- Error: { "error": "<message>" }
- 401 `Membutuhkan autentikasi` — missing or malformed Authorization header
- 401 `Token tidak valid atau kadaluwarsa` — invalid/expired JWT
- 403 `Anda tidak memiliki izin akses` — role not authorized for the endpoint

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
- Cache-Control: no-cache
- Connection: keep-alive

Event payload (example):
```json
{
  "action": "UPDATE_HEWAN",
  "data": { }
}
```

Actions:
- UPDATE_HEWAN — sent when any hewan mutation occurs (progress, timbang, packing, kelengkapan)
- UPDATE_DASHBOARD — sent after hewan/distribusi mutations (recalculated summary)
- UPDATE_DISTRIBUSI — sent when distribusi is updated

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

Notes:
- `total_hewan_selesai` counts hewan where `waktu_selesai_kuliti IS NOT NULL`
- `waktu_mulai` is `MIN(waktu_mulai_jagal)` across all hewan
- `waktu_selesai` is `MAX(waktu_selesai_kuliti)` — only set when all hewan are complete

### GET /dashboard/hewan
Public list of hewan (includes pengawas).

Query params:
- search: string (matches kode_hewan, nama_sohibul, catatan)
- tipe: qurban | sedekah
- jenis_hewan: sapi | kambing
- pengawas_id: number (exact ID match; if not a valid number, treated as nama_lengkap search)
- pengawas: string (matches pengawas nama_lengkap with LIKE)

Ordering:
1. In-progress hewan first (waktu_mulai_jagal set, waktu_selesai_kuliti not set)
2. Not started hewan second (waktu_mulai_jagal not set)
3. Completed hewan last
4. Then by waktu_mulai_jagal ASC, kode_hewan ASC

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

Errors:
- 500 Gagal mengambil data hewan

## Authenticated Endpoints

### GET /hewan
Requires any role: admin, pengawas, jagal, kulit, cacah_daging, cacah_tulang, packing, distribusi.

> Note: `koordinator_pengawas` is **NOT** included in the authorized roles for this endpoint.

Query params:
- search: string (matches kode_hewan, nama_sohibul, catatan)
- tipe: qurban | sedekah
- jenis_hewan: sapi | kambing
- pengawas_id: number (ignored for role pengawas, forced to own ID)

Ordering (role-based):
- Role **jagal** / **admin** / **pengawas** / **koordinator_pengawas**: sorted by jagal timestamps
- Role **kulit**: sorted by kuliti timestamps
- Role **cacah_daging**: sorted by cacah_daging timestamps
- Role **cacah_tulang**: sorted by cacah_tulang timestamps
- Role **packing** / **distribusi**: sorted by packing timestamps
- Within each: in-progress first → not started → completed, then by start time ASC, kode_hewan ASC

Response 200: same shape as GET /dashboard/hewan

Errors:
- 500 Gagal mengambil data hewan

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
- Cannot start a pos that is already started
- Cannot finish a pos that hasn't started or is already finished

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
- 400 Proses jagal sudah dimulai
- 400 Status jagal tidak valid untuk diselesaikan
- 400 Proses kulit sudah dimulai
- 400 Status kulit tidak valid untuk diselesaikan
- 400 Proses cacah daging sudah dimulai
- 400 Status cacah daging tidak valid untuk diselesaikan
- 400 Proses cacah tulang sudah dimulai
- 400 Status cacah tulang tidak valid untuk diselesaikan
- 400 Pos operasional tidak dikenali
- 403 Anda tidak berhak mengelola hewan ini
- 404 Hewan tidak ditemukan
- 500 Gagal menyimpan progress

### PATCH /hewan/:id/timbang
Update berat_daging and berat_tulang. Jagal must be completed.

Rules:
- Jagal must be completed
- For role pengawas: can only update own hewan

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

Errors:
- 400 Kirimkan berat_daging dan berat_tulang berupa angka
- 400 Proses jagal belum selesai
- 403 Anda tidak berhak mengelola hewan ini
- 404 Hewan tidak ditemukan
- 500 Gagal menyimpan data timbangan

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

Errors:
- 400 Payload tidak valid
- 403 Anda tidak berhak mengelola hewan ini
- 404 Hewan tidak ditemukan
- 500 Gagal menyimpan kelengkapan

## Packing Endpoints (admin, packing)

### PATCH /hewan/:id/packing
Start or finish packing.

Start packing:
```json
{ "status": "mulai" }
```

Finish packing:
```json
{ "status": "selesai", "total_kantong": 12 }
```

Rules:
- Jagal must be completed
- Cannot start if packing already started
- Cannot finish if packing hasn't started
- total_kantong required when status=selesai
- If packing was already finished, total_kantong is updated but waktu_selesai_packing stays

Response 200:
```json
{
  "message": "Data packing berhasil diperbarui",
  "data": { }
}
```

Errors:
- 400 Payload tidak valid
- 400 Proses jagal belum selesai
- 400 Proses packing sudah dimulai
- 400 Proses packing belum dimulai
- 400 Total kantong wajib diisi saat menyelesaikan packing
- 404 Hewan tidak ditemukan
- 500 Gagal menyimpan data packing

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

Errors:
- 400 Gagal mengambil data distribusi

### GET /distribusi/:user_id
Get distribusi for a user. If no data, returns placeholder with jumlah_kantong=0.

Response 200 (exists):
```json
{
  "message": "Berhasil",
  "data": {
    "id": 1,
    "user_id": 5,
    "jumlah_kantong": 10,
    "user": { },
    "created_at": "2026-05-08T10:00:00Z",
    "updated_at": "2026-05-08T10:00:00Z"
  }
}
```

Response 200 (no record):
```json
{
  "message": "Belum ada distribusi",
  "data": { "JumlahKantong": 0, "UserID": "5" }
}
```

Errors:
- 400 Gagal mengambil data distribusi

### PATCH /distribusi/:user_id
Increment or decrement jumlah_kantong. Creates record if it doesn't exist (FirstOrCreate).

Request body:
```json
{ "penambahan": 3 }
```

Rules:
- Role distribusi can only update own user_id
- jumlah_kantong cannot go below 0 (clamped to 0)

Response 200:
```json
{
  "message": "Data distribusi berhasil diupdate",
  "data": { }
}
```

Errors:
- 400 Kirimkan parameter 'penambahan'
- 400 ID User tidak valid
- 403 Anda hanya dapat mengupdate data distribusi Anda sendiri
- 500 Gagal mengakses data distribusi
- 500 Gagal menyimpan data distribusi

## Admin Endpoints (admin only)

### GET /users
Query params:
- role: role string
- search: username substring (LIKE match)

Response 200:
```json
{ "message": "Berhasil", "data": [ ] }
```

Errors:
- 500 Gagal mengambil data user

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

Errors:
- 400 Input tidak valid
- 400 Role tidak valid
- 500 Gagal membuat user baru

### PUT /users/:id
Request body:
```json
{
  "username": "string",
  "password": "string",
  "role": "admin"
}
```

> Note: `nama_lengkap` is **NOT** updatable via this endpoint. `password` is optional — if empty string or omitted, the existing password is kept.

Response 200:
```json
{ "message": "User berhasil diperbarui" }
```

Errors:
- 400 Input tidak valid
- 404 User tidak ditemukan
- 500 Gagal mengupdate user

### DELETE /users/:id
Response 200:
```json
{ "message": "User berhasil dihapus" }
```

Errors:
- 404 User tidak ditemukan
- 500 Gagal menghapus user

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

Validation:
- kode_hewan: required
- tipe: required, must be `qurban` or `sedekah`
- jenis_hewan: required, must be `sapi` or `kambing`
- nama_sohibul: required
- pengawas_id: required

Response 201:
```json
{ "message": "Data hewan berhasil diinput", "data": { } }
```

Errors:
- 400 Input tidak valid
- 500 Gagal menyimpan data hewan

### PUT /hewan/:id

Request body (all fields applied directly, not partial-update):
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

Errors:
- 400 Input tidak valid
- 404 Data hewan tidak ditemukan
- 500 Gagal mengupdate data hewan

### DELETE /hewan/:id
Rule: cannot delete if waktu_mulai_jagal is set.

Response 200:
```json
{ "message": "Data hewan berhasil dihapus" }
```

Errors:
- 403 Hewan sudah diproses, data tidak dapat dihapus
- 404 Data hewan tidak ditemukan
- 500 Gagal menghapus data hewan
