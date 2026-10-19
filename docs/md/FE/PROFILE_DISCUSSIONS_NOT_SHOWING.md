# FE Issue: Diskusi Tidak Muncul di Profile

## Status Saat Ini

Endpoint profile saat ini hanya mengembalikan data user dan statistik. Data diskusi milik user belum dikembalikan oleh endpoint profile.

Endpoint yang tersedia:

```http
GET /api/v1/users/me
GET /api/v1/users/me/stats
```

`GET /api/v1/users/me/stats` hanya mengembalikan jumlah diskusi melalui field:

```json
{
  "discussion_count": 0
}
```

Field tersebut adalah jumlah, bukan daftar diskusi.

## Penyebab Diskusi Tidak Tampil

Endpoint daftar diskusi saat ini adalah:

```http
GET /api/v1/discussions
```

Endpoint tersebut mendukung filter:

- `search`
- `category`
- `role`
- `sort`
- `page`
- `limit`

Belum tersedia filter `user_id` atau endpoint khusus seperti:

```http
GET /api/v1/users/me/discussions
```

Karena itu FE tidak dapat mengambil daftar diskusi milik user secara spesifik dengan kontrak API saat ini.

## Yang Dapat Dilakukan FE Saat Ini

Gunakan `discussion_count` untuk menampilkan jumlah diskusi pada profile:

```javascript
const response = await api.get("/users/me/stats");
const discussionCount = response.data.data.discussion_count;
```

Untuk daftar diskusi, jangan mengambil seluruh `/discussions` lalu melakukan filter berdasarkan creator di FE. Cara tersebut tidak aman dan tidak menjamin semua data user tersedia karena endpoint menggunakan pagination.

## Kebutuhan Backend

Agar diskusi dapat ditampilkan di profile, backend perlu menyediakan salah satu endpoint berikut:

### Profile user yang sedang login

```http
GET /api/v1/users/me/discussions?page=1&limit=10
Authorization: Bearer ACCESS_TOKEN
```

### Profile user berdasarkan ID

```http
GET /api/v1/users/:id/discussions?page=1&limit=10
```

Response yang disarankan:

```json
{
  "data": {
    "items": [],
    "page": 1,
    "limit": 10,
    "total": 0,
    "total_pages": 0
  },
  "message": "User discussions retrieved successfully"
}
```

## Checklist FE

- Pastikan request profile menggunakan access token terbaru.
- Gunakan `discussion_count` hanya untuk angka statistik.
- Jangan menganggap `GET /users/me` berisi daftar diskusi.
- Tampilkan empty state jika daftar diskusi kosong.
- Setelah backend menyediakan endpoint khusus, tambahkan request terpisah untuk daftar diskusi profile.

## Troubleshooting: Diskusi Tercatat sebagai Milik User Lain

Backend menentukan pemilik diskusi dari `user_id` yang terdapat pada access token. FE tidak perlu dan tidak boleh mengirim `user_id` pada body create discussion.

### Header yang benar

```http
Authorization: Bearer ACCESS_TOKEN
```

Gunakan nilai:

```text
data.access_token
```

Jangan menggunakan `refresh_token` sebagai Bearer token.

### Validasi user yang sedang login

Sebelum membuat diskusi, cek identitas token:

```http
GET /api/v1/users/me
Authorization: Bearer ACCESS_TOKEN
```

Email dan ID dari response harus sesuai dengan akun yang sedang aktif di FE.

### Create discussion

```http
POST /api/v1/discussions
Authorization: Bearer ACCESS_TOKEN
Content-Type: application/json
```

Body hanya berisi data diskusi:

```json
{
  "title": "Judul diskusi",
  "description": "Isi diskusi",
  "category": "technology",
  "tags": "go, backend"
}
```

Jangan mengirim field berikut dari FE:

```json
{
  "user_id": "...",
  "author_id": "...",
  "creator_id": "..."
}
```

### Penyebab umum

- FE masih menggunakan access token dari akun sebelumnya.
- Token tersimpan di local storage belum dibersihkan saat logout.
- FE memakai `refresh_token` sebagai access token.
- Environment Bruno menggunakan token user lain.
- Profile menampilkan data creator dari object user yang salah.

Jika diskusi tetap tercatat sebagai user lain, bandingkan response `GET /users/me` dengan field `author_id` pada response diskusi.
